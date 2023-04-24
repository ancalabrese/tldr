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

type Embedding openai.Embedding

type RelatedEmbedding struct {
	Embedding      Embedding
	RelevanceScore float32
	Query          string
}

func GetKbEmbeddings(r io.Reader) ([]Embedding, error) {
	kbEmbeddings := make([]Embedding, 0)
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

		kbEmbeddings = append(kbEmbeddings, Embedding{
			Object:    entry[0],
			Embedding: convertedEmbeddings,
		})
	}
	return kbEmbeddings, err
}

func GetQueryEmbedding(ctx context.Context, q string, llm openai.Client) (Embedding, error) {
	queryEmbeddingReq := openai.EmbeddingRequest{
		Input: []string{q},
		Model: openai.AdaEmbeddingV2,
	}

	queryEmbeddingRes, err := llm.CreateEmbeddings(ctx, queryEmbeddingReq)
	if err != nil {
		return Embedding{}, err
	}

	return Embedding{
		Object:    queryEmbeddingRes.Object,
		Embedding: queryEmbeddingRes.Data[0].Embedding,
	}, nil
}

func RankEmbeddingsByRelatedness(
	llm openai.Client,
	kbEmbeddings []Embedding,
	queryEmbedding Embedding,
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
