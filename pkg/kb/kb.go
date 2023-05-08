package kb

import (
	"io"
	"net/url"
	"os"

	"github.com/sashabaranov/go-openai"
)

type Kb struct {
	uri               url.URL
	Embeddings        []openai.Embedding
	embeddingFilePath string
}

func New(uri url.URL) *Kb {
	return &Kb{
		uri: uri,
	}
}

func (kb *Kb) GetKbReader() (io.ReadCloser, error) {
	fp := kb.uri.Path
	return os.Open(fp)
}

func (kb *Kb) GetKb() (io.ReadWriteCloser, error) {
	fp := kb.uri.Path
	return os.OpenFile(fp, os.O_RDWR, 0)
}
