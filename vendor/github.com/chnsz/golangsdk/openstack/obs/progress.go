package obs

import (
	"io"
)

type ProgressEventType int

type ProgressEvent struct {
	ConsumedBytes int64
	TotalBytes    int64
	EventType     ProgressEventType
}

const (
	TransferStartedEvent ProgressEventType = 1 + iota
	TransferDataEvent
	TransferCompletedEvent
	TransferFailedEvent
)

func newProgressEvent(eventType ProgressEventType, consumed, total int64) *ProgressEvent {
	return &ProgressEvent{
		ConsumedBytes: consumed,
		TotalBytes:    total,
		EventType:     eventType,
	}
}

type ProgressListener interface {
	ProgressChanged(event *ProgressEvent)
}

type readerTracker struct {
	completedBytes int64
}

// publishProgress
func publishProgress(listener ProgressListener, event *ProgressEvent) {
	if listener != nil && event != nil {
		listener.ProgressChanged(event)
	}
}

type teeReader struct {
	reader        io.Reader
	consumedBytes int64
	totalBytes    int64
	tracker       *readerTracker
	listener      ProgressListener
}

func TeeReader(reader io.Reader, totalBytes int64, listener ProgressListener, tracker *readerTracker) io.ReadCloser {
	return &teeReader{
		reader:        reader,
		consumedBytes: 0,
		totalBytes:    totalBytes,
		tracker:       tracker,
		listener:      listener,
	}
}

func (t *teeReader) Read(p []byte) (n int, err error) {
	n, err = t.reader.Read(p)

	if err != nil && err != io.EOF {
		event := newProgressEvent(TransferFailedEvent, t.consumedBytes, t.totalBytes)
		publishProgress(t.listener, event)
	}

	if n > 0 {
		t.consumedBytes += int64(n)

		if t.listener != nil {
			event := newProgressEvent(TransferDataEvent, t.consumedBytes, t.totalBytes)
			publishProgress(t.listener, event)
		}

		if t.tracker != nil {
			t.tracker.completedBytes = t.consumedBytes
		}
	}

	return
}

func (r *teeReader) Size() int64 {
	return r.totalBytes
}

func (t *teeReader) Close() error {
	if rc, ok := t.reader.(io.ReadCloser); ok {
		return rc.Close()
	}
	return nil
}
