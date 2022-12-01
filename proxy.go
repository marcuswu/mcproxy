package main

import (
	"errors"
	"os"
	"sync"

	"github.com/marcuswu/mcproxy/cheat"
	_ "github.com/marcuswu/mcproxy/world"
	"github.com/marcuswu/mcproxy/world/chunk"
	"github.com/marcuswu/mcproxy/token"
	"github.com/pelletier/go-toml"
	"github.com/rs/zerolog/log"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"golang.org/x/oauth2"
)

func main() {
	config := readConfig()
	src := token.TokenSource()

	p, err := minecraft.NewForeignStatusProvider(config.Connection.RemoteAddress)
	if err != nil {
		panic(err)
	}
	listener, err := minecraft.ListenConfig{
		StatusProvider: p,
	}.Listen("raknet", config.Connection.LocalAddress)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	for {
		c, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go handleConn(c.(*minecraft.Conn), listener, config, src)
	}
}

func handleConn(clientConn *minecraft.Conn, listener *minecraft.Listener, config config, src oauth2.TokenSource) {
	serverConn, err := minecraft.Dialer{
		TokenSource: src,
		ClientData:  clientConn.ClientData(),
	}.Dial("raknet", config.Connection.RemoteAddress)
	if err != nil {
		panic(err)
	}
	var g sync.WaitGroup
	g.Add(2)
	go func() {
		if err := clientConn.StartGame(serverConn.GameData()); err != nil {
			panic(err)
		}
		g.Done()
	}()
	go func() {
		if err := serverConn.DoSpawn(); err != nil {
			panic(err)
		}
		g.Done()
	}()

	proxy := &cheat.Proxy{ClientConn: clientConn, ServerConn: serverConn, Chunks: make(map[protocol.ChunkPos]*chunk.Chunk)}

	go func() {
		defer listener.Disconnect(clientConn, "connection lost")
		defer serverConn.Close()
		for {
			pk, err := clientConn.ReadPacket()
			if err != nil {
				return
			}
			pk, forward := proxy.HandleClientPacket(pk)
			if !forward {
				continue
			}
			if err := serverConn.WritePacket(pk); err != nil {
				if disconnect, ok := errors.Unwrap(err).(minecraft.DisconnectError); ok {
					_ = listener.Disconnect(clientConn, disconnect.Error())
				}
				return
			}
		}
	}()
	go func() {
		defer serverConn.Close()
		defer listener.Disconnect(clientConn, "connection lost")
		for {
			pk, err := serverConn.ReadPacket()
			if err != nil {
				if disconnect, ok := errors.Unwrap(err).(minecraft.DisconnectError); ok {
					_ = listener.Disconnect(clientConn, disconnect.Error())
				}
				return
			}
			pk, forward := proxy.HandleServerPacket(pk)
			if !forward {
				continue
			}
			if err := clientConn.WritePacket(pk); err != nil {
				return
			}
		}
	}()
}

type config struct {
	Connection struct {
		LocalAddress  string
		RemoteAddress string
	}
}

func readConfig() config {
	c := config{}
	if _, err := os.Stat("config.toml"); os.IsNotExist(err) {
		f, err := os.Create("config.toml")
		if err != nil {
			log.Error().Msgf("error creating config: %v", err)
			os.Exit(1)
		}
		data, err := toml.Marshal(c)
		if err != nil {
			log.Error().Msgf("error encoding default config: %v", err)
			os.Exit(1)
		}
		if _, err := f.Write(data); err != nil {
			log.Error().Msgf("error writing encoded default config: %v", err)
			os.Exit(1)
		}
		_ = f.Close()
	}
	data, err := os.ReadFile("config.toml")
	if err != nil {
		log.Error().Msgf("error reading config: %v", err)
		os.Exit(1)
	}
	if err := toml.Unmarshal(data, &c); err != nil {
		log.Error().Msgf("error decoding config: %v", err)
		os.Exit(1)
	}
	if c.Connection.LocalAddress == "" {
		c.Connection.LocalAddress = "0.0.0.0:19132"
	}
	data, _ = toml.Marshal(c)
	if err := os.WriteFile("config.toml", data, 0644); err != nil {
		log.Error().Msgf("error writing config file: %v", err)
	}
	return c
}
