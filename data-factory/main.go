// package datafactory

// import "fmt"

// type Factory[K any] struct {
// 	name          string
// 	lines         []string
// 	data          []map[string]interface{}
// 	lineProcessor func(a []K) ([]string, error)
// }

// type DataFactory[T any] func() *Factory[T]

// func GetInstance[T any](name string) DataFactory[T] {
// 	var instance *Factory[T]
// 	return func() *Factory[T] {
// 		if instance != nil {
// 			fmt.Println("returning cached")
// 			return instance
// 		}
// 		fmt.Println("instantiating")
// 		instance = &Factory[T]{
// 			name: name,
// 		}
// 		return instance
// 	}
// }

// // func (s *DataFactory[T]) GetName() string {
// // 	return *s. .name
// // }

// // func (s *DataFactory) GetLines() []string {
// // 	return s.lines
// // }

// // func (s *DataFactory) SetLineProcessor(lineProcessor func(a) ([]string, error)) {
// // 	s.lineProcessor = lineProcessor
// // }

// // func (s *DataFactory) LoadData(data []map[string]interface{}) {
// // 	s.data = data
// // }
