package performance

import "testing"

// func BenchmarkRCache(b *testing.B) {
// 	r := RCache{}
// 	r.Set("value")
// 	for n := 0; n < b.N; n++ {
// 		go r.Set("BBB")
// 		go r.Get("value")
// 	}
// }

// func BenchmarkLCache(b *testing.B) {
// 	r := LCache{}
// 	r.Set("value")
// 	for n := 0; n < b.N; n++ {
// 		go r.Set("BBB")
// 		go r.Get("value")
// 	}
// }

func BenchmarkLocalCache(b *testing.B) {
	var LocalCache LocalModelStruct = LocalModelStruct{LocalData: make(map[string]LocalDataStruct)}
	LocalCache.SetValue("somekey" , "value" , 1)
	for n := 0; n < b.N; n++ {
		go LocalCache.SetValue("somekey" , "value" , 1)
		go LocalCache.GetValue("somekey")
	}
}
