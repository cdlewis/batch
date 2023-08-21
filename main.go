package main

import (
	"fmt"
	"sync"
)

func main() {
	dataSource := NewFanoutDataSource(
		[]DataSource[int, []int]{
			TransformDataSourceResult[int, []int](
				TransformDataSourceResult[int, []int](
					NewFruitDataSource([]int{1, 2, 3}),
					func(result []int) []int {
						for idx := range result {
							result[idx] /= 2
						}
						return result
					},
				),
				func(result []int) []int {
					return result[1:]
				},
			),
			TransformDataSourceResult[int, []int](
				NewFruitDataSource([]int{2, 4, 6}),
				func(result []int) []int {
					for idx := range result {
						result[idx] *= 2
					}
					return result
				},
			),
		},
	)

	fmt.Println("Initial Graph:\n")
	DependencyGraph(dataSource, 0)
	//fmt.Println(dataSource.Fetch(1))
}

func DependencyGraph(ds WrappedDataSource, depth int) {
	for i := 0; i < depth; i++ {
		fmt.Print("\t")
	}
	fmt.Println(ds.Type())

	if wrappedDS, ok := ds.(WrappedDataSource); ok {
		for _, child := range wrappedDS.Unwrap() {
			DependencyGraph(child, depth+1)
		}
	}
}

type FanoutDataSource[Q, R any] struct {
	DataSource[Q, []R]
	WrappedDataSource

	resultTransformers []func(Q, any) R
	dataSources        []DataSource[Q, R]
}

func NewFanoutDataSource[Q, R any](ds []DataSource[Q, R]) DataSource[Q, []R] {
	return FanoutDataSource[Q, R]{
		dataSources: ds,
	}
}

func (f FanoutDataSource[Q, R]) Type() string {
	return "FanoutDataSource"
}

func (f FanoutDataSource[Q, R]) Fetch(q Q) []R {
	wg := sync.WaitGroup{}
	results := make([]R, len(f.dataSources))

	for idx, f := range f.dataSources {
		idx := idx
		f := f
		wg.Add(1)

		go func() {
			results[idx] = f.Fetch(q)
			wg.Done()
		}()
	}

	wg.Wait()

	return results
}

func (f FanoutDataSource[Q, R]) Unwrap() []WrappedDataSource {
	children := make([]WrappedDataSource, 0, len(f.dataSources))
	for _, d := range f.dataSources {
		children = append(children, d)
	}
	return children
}

type WrappedDataSource interface {
	Type() string
	Unwrap() []WrappedDataSource
}

type DataSource[Q, R any] interface {
	Type() string

	Fetch(Q) R

	Unwrap() []WrappedDataSource
}

type BatchDataSource[Q, R any] interface {
	DataSource[Q, R]

	AppendQuery(id int, current any, base any) any

	FetchBatch(any) map[int]R
}

type AppleDataSource struct {
	BatchDataSource[int, []int]

	numbersToReturn []int

	results chan []int
}

func NewFruitDataSource(numbers []int) AppleDataSource {
	return AppleDataSource{numbersToReturn: numbers}
}

func (a AppleDataSource) Type() string {
	return "Apple"
}

func (a AppleDataSource) Fetch(q int) []int {
	return <-a.results
}

func (a AppleDataSource) AppendQuery(id int, current map[int][]int) (map[int][]int, chan []int) {
	if current == nil {
		return map[int][]int{
			id: a.numbersToReturn,
		}, a.results
	}
	current[id] = a.numbersToReturn
	return current, a.results
}

func (a AppleDataSource) FetchBatch(q map[int][]int) map[int][]int {
	return q
}

func (a AppleDataSource) Unwrap() []WrappedDataSource {
	return []WrappedDataSource{}
}

type StackedDataSource[Q1, R1, R2 any] struct {
	DataSource[Q1, R2]

	ds1 DataSource[Q1, R1]

	ds2 DataSource[R1, R2]
}

func NewStackedDataSource[Q1, R1, R2 any](ds1 DataSource[Q1, R1], ds2 DataSource[R1, R2]) DataSource[Q1, R2] {
	return StackedDataSource[Q1, R1, R2]{
		ds1: ds1,
		ds2: ds2,
	}
}

func (s StackedDataSource[Q1, R1, R2]) Type() string {
	return "StackedDataSource"
}

func (s StackedDataSource[Q1, R1, R2]) Fetch(q1 Q1) R2 {
	r1 := s.ds1.Fetch(q1)
	return s.ds2.Fetch(r1)
}

func (s StackedDataSource[Q1, R1, R2]) Unwrap() []WrappedDataSource {
	return []WrappedDataSource{s.ds1, s.ds2}
}

type TransformResultDataSource[Q, R1, R2 any] struct {
	DataSource[Q, R2]

	baseDataSource DataSource[Q, R1]

	transformer func(r1 R1) R2
}

func TransformDataSourceResult[Q, R1, R2 any](ds DataSource[Q, R1], fn func(R1) R2) DataSource[Q, R2] {
	return TransformResultDataSource[Q, R1, R2]{
		baseDataSource: ds,
		transformer:    fn,
	}
}

func (t TransformResultDataSource[Q, R1, R2]) Type() string {
	return "TransformResult"
}

func (t TransformResultDataSource[Q, R1, R2]) Fetch(q Q) R2 {
	return t.transformer(t.baseDataSource.Fetch(q))
}

func (t TransformResultDataSource[Q, R1, R2]) Unwrap() []WrappedDataSource {
	return []WrappedDataSource{t.baseDataSource}
}

type Orange struct{}

type OrangeDataSource struct{}

func (o OrangeDataSource) Type() string {
	return "Orange"
}

func (o OrangeDataSource) Fetch(q int) []Orange {
	result := []Orange{}
	for i := 0; i < q; i++ {
		result = append(result, Orange{})
	}
	return result
}
