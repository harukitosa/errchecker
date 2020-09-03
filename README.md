# errchecker

## 機能

- error が返り値として存在する関数のうち nil でしか返していないものを指摘します
- Pointing out a function with a return value of error that only returns nil

## sample

対象ファイル

```go
package sample_test1

func sample1() error {
	if false {
		s := 0
		s++
		if s >= 1 {
			s++
		} else {
			s++
			for {
				return nil
			}
		}
	}
	return nil
}
```

実行結果

```zsh
go vet -vettool=$(which errchecker) testdata/src/sample_test/sample_test1.go
# command-line-arguments
testdata/src/sample_test/sample_test1.go:3:1: It returns nil in all the places where it should return error. Please fix the return value
```
