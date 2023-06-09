package data

import (
	"context"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/ancalabrese/tldr/codecs"
	"github.com/sashabaranov/go-openai"
)

type DistanceFunc func(embeddingsA []float32, embeddingsB []float32) float32

type RelatedEmbedding struct {
	Embedding      openai.Embedding
	RelevanceScore float32
	Query          string
}

func ParseEmbeddings(r io.Reader) ([]openai.Embedding, error) {
	kbEmbeddings := make([]openai.Embedding, 0)
	kb, err := codecs.CsvReaderFunc(r)
	if err != nil {
		return nil, err
	}

	for _, entry := range kb {
		embeddings := strings.Split(entry[1], ",")
		convertedEmbeddings := make([]float32, 0)

		for _, val := range embeddings {
			convertedVal, err := strconv.ParseFloat(val, 32)
			if err != nil {
				continue
			}
			convertedEmbeddings = append(convertedEmbeddings, float32(convertedVal))
		}

		kbEmbeddings = append(kbEmbeddings, openai.Embedding{
			Object:    entry[0],
			Embedding: convertedEmbeddings,
		})
	}
	return kbEmbeddings, err
}

func GetEmbeddings(ctx context.Context, text []string, llm *openai.Client) ([]openai.Embedding, error) {
	queryEmbeddingReq := openai.EmbeddingRequest{
		Input: text,
		Model: openai.AdaEmbeddingV2,
	}

	queryEmbeddingRes, err := llm.CreateEmbeddings(ctx, queryEmbeddingReq)
	if err != nil {
		return nil, err
	}
	return queryEmbeddingRes.Data, nil
}

func RankEmbeddingsByRelatedness(
	llm openai.Client,
	kbEmbeddings []openai.Embedding,
	queryEmbedding openai.Embedding,
	distanceFunc DistanceFunc,
	maxResults int,
	ctx context.Context) []RelatedEmbedding {

	relatedEmbeddings := make([]RelatedEmbedding, 0)

	for _, kbe := range kbEmbeddings {
		score := distanceFunc(kbe.Embedding, queryEmbedding.Embedding)
		relatedEmbedding := RelatedEmbedding{
			Query:          queryEmbedding.Object,
			Embedding:      kbe,
			RelevanceScore: score,
		}
		relatedEmbeddings = append(relatedEmbeddings, relatedEmbedding)
	}

	sort.Slice(relatedEmbeddings, func(i, j int) bool {
		return relatedEmbeddings[i].RelevanceScore > relatedEmbeddings[j].RelevanceScore
	})

	return relatedEmbeddings[:maxResults]
}
