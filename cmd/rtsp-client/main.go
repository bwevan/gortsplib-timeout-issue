package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/bluenviron/gortsplib/v5"
	"github.com/bluenviron/gortsplib/v5/pkg/base"
	"github.com/bluenviron/gortsplib/v5/pkg/format"
	"github.com/pion/rtp"
)

func run() error {
	var (
		rtspURL     string
		rtspTimeout time.Duration
	)
	flag.StringVar(&rtspURL, "rtsp-url", "rtsp://localhost:8554/stream", "RTSP URL to connect to")
	flag.DurationVar(&rtspTimeout, "rtsp-timeout", 10*time.Second, "RTSP timeout")
	flag.Parse()

	u, err := base.ParseURL(rtspURL)
	if err != nil {
		return fmt.Errorf("parse RTSP URL: %w", err)
	}

	log.Printf("Dialing RTSP URL %s with timeout %s", u.String(), rtspTimeout.String())

	protocol := gortsplib.ProtocolTCP
	c := gortsplib.Client{
		Scheme:      u.Scheme,
		Host:        u.Host,
		Protocol:    &protocol,
		ReadTimeout: rtspTimeout,
	}

	if err = c.Start(); err != nil {
		return fmt.Errorf("RTSP start: %w", err)
	}
	defer c.Close()

	desc, _, err := c.Describe(u)
	if err != nil {
		return fmt.Errorf("RTSP describe: %w", err)
	}

	var f *format.Generic
	media := desc.FindFormat(&f)
	if media == nil {
		return errors.New("no generic format found in media description")
	}

	if _, err = c.Setup(u, media, 0, 0); err != nil {
		return fmt.Errorf("RTSP setup: %w", err)
	}

	c.OnPacketRTP(media, f, func(p *rtp.Packet) {
		log.Printf("Received RTP packet with timestamp %d, payload %s\n", p.Timestamp, string(p.Payload))
	})

	if _, err = c.Play(nil); err != nil {
		return fmt.Errorf("RTSP play: %w", err)
	}

	if err = c.Wait(); err != nil {
		return fmt.Errorf("RTSP wait: %w", err)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("Failed to run RTSP client: %s", err)
	}
}
