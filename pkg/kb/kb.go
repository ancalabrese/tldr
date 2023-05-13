package kb

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"

	"github.com/ancalabrese/tldr/data"
	"github.com/sashabaranov/go-openai"
)

type Kb struct {
	uri        *url.URL
	Embeddings []openai.Embedding
}

func New(uri *url.URL) *Kb {
	kb := &Kb{
		uri: uri,
	}
	return kb
}

func (kb *Kb) Parse(ctx context.Context, llm *openai.Client) error {
	content, err := kb.parseContent()
	if err != nil {
		return err
	}

	emb, err := data.GetEmbeddings(ctx, content, llm)
	if err != nil {
		return fmt.Errorf("couldn't get kb embeddings: %w", err)
	}

	kb.Embeddings = emb
	return nil
}

func (kb *Kb) GetKbReader() (io.ReadCloser, error) {
	fp := kb.uri.Path
	return os.Open(fp)
}

func (kb *Kb) GetKb() (io.ReadWriteCloser, error) {
	fp := kb.uri.Path
	return os.OpenFile(fp, os.O_RDWR, 0)
}

func (kb *Kb) parseContent() ([]string, error) {

	r, err := kb.GetKbReader()
	if err != nil {
		return nil, fmt.Errorf("kb content parsing failed: %w", err)
	}
	defer r.Close()

	content := make([]string, 0)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		content = append(content, scanner.Text())
	}

	// If there was an error and we didn't parse anything return the error.
	// Otherwise work with what we have
	if scanner.Err() != nil && len(content) == 0 {
		return nil, fmt.Errorf("failed to parse kb: %w", err)
	} else if scanner.Err() != nil && len(content) > 0 {
		log.Println("[Warning] - KB partially parsed. Err:", err)
	}

	return content, nil
}
