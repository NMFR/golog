package repositories // import github.com/mlimaloureiro/golog/repositories

import "io"

func trySeek(readWriter io.ReadWriter, offset int64, whence int) (int64, error) {
	seeker, isSeekable := readWriter.(io.Seeker)
	if isSeekable == false {
		return 0, nil
	}
	return seeker.Seek(offset, whence)
}
