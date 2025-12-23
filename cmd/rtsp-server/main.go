package main

import (
	"log"
	"sync"
	"time"

	"github.com/bluenviron/gortsplib/v5"
	"github.com/bluenviron/gortsplib/v5/pkg/base"
	"github.com/bluenviron/gortsplib/v5/pkg/description"
	"github.com/bluenviron/gortsplib/v5/pkg/format"
	"github.com/pion/rtp"
)

func main() {
	h := &serverHandler{}

	h.mutex.Lock()

	h.server = &gortsplib.Server{
		Handler:        h,
		RTSPAddress:    ":8554",
		UDPRTPAddress:  ":8000",
		UDPRTCPAddress: ":8001",
		// Disable RTCP sender reports so time-out is triggered
		DisableRTCPSenderReports: true,
	}

	if err := h.server.Start(); err != nil {
		panic(err)
	}
	defer h.server.Close()

	desc := &description.Session{
		Medias: []*description.Media{{
			Type: description.MediaTypeApplication,
			Formats: []format.Format{&format.Generic{
				RTPMa:      "vnd.onvif.metadata/90000",
				PayloadTyp: 98,
			}},
		}},
	}

	h.stream = &gortsplib.ServerStream{
		Server: h.server,
		Desc:   desc,
	}

	if err := h.stream.Initialize(); err != nil {
		panic(err)
	}
	defer h.stream.Close()

	go streamMockMetadataPayload(desc.Medias[0], h.stream)

	h.mutex.Unlock()

	log.Printf("server is ready on %s", h.server.RTSPAddress)
	panic(h.server.Wait())
}

// serverHandler is a very simple handler implementation for *gortsplib.Server, usable for testing.
type serverHandler struct {
	server *gortsplib.Server
	stream *gortsplib.ServerStream
	mutex  sync.RWMutex
}

// OnDescribe handles RTSP DESCRIBE requests.
func (h *serverHandler) OnDescribe(_ *gortsplib.ServerHandlerOnDescribeCtx) (
	*base.Response,
	*gortsplib.ServerStream,
	error,
) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	return &base.Response{
		StatusCode: base.StatusOK,
	}, h.stream, nil
}

// OnSetup handles RTSP SETUP requests.
func (h *serverHandler) OnSetup(_ *gortsplib.ServerHandlerOnSetupCtx) (*base.Response, *gortsplib.ServerStream, error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	return &base.Response{
		StatusCode: base.StatusOK,
	}, h.stream, nil
}

// OnPlay handles RTSP PLAY requests.
func (h *serverHandler) OnPlay(_ *gortsplib.ServerHandlerOnPlayCtx) (*base.Response, error) {
	return &base.Response{
		StatusCode: base.StatusOK,
	}, nil
}

func streamMockMetadataPayload(m *description.Media, s *gortsplib.ServerStream) {
	sequence := uint16(1)
	payload := []byte("foobar")

	for {
		p := &rtp.Packet{
			Header: rtp.Header{
				Version:        2,
				PayloadType:    98,
				SequenceNumber: sequence,
				SSRC:           123,
				Marker:         true,
			},
			Payload: payload,
		}

		if err := s.WritePacketRTP(m, p); err != nil {
			panic(err)
		}

		time.Sleep(5 * time.Second)
		sequence++
	}
}
