package mocka

import (
	"os"
	"reflect"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("mockFileInfo", func() {
	var fileInfo mockFileInfo

	BeforeEach(func() {
		fileInfo = mockFileInfo{
			name: "i_am_a_file",
			size: int64(0),
		}
	})

	It("implements os.FileInfo", func() {
		mockFileInfoType := reflect.TypeOf(&fileInfo)
		fileInfoType := reflect.TypeOf((*os.FileInfo)(nil)).Elem()
		Expect(mockFileInfoType.Implements(fileInfoType)).To(BeTrue())
	})

	It("returns the name of the mockFileInfo", func() {
		Expect(fileInfo.Name()).To(Equal("i_am_a_file"))
		fileInfo.name = "i_am_a_different_file"
		Expect(fileInfo.Name()).To(Equal("i_am_a_different_file"))
	})

	It("returns the size of the mockFileInfo", func() {
		Expect(fileInfo.Size()).To(Equal(int64(0)))
		fileInfo.size = int64(100)
		Expect(fileInfo.Size()).To(Equal(int64(100)))
	})

	It("should return a regular file FileMode", func() {
		fileMode := fileInfo.Mode()

		Expect(fileMode.IsRegular()).To(BeTrue())
	})

	It("returns false for .isDir()", func() {
		Expect(fileInfo.IsDir()).To(BeFalse())
	})

	It("returns nil for .Sys()", func() {
		Expect(fileInfo.Sys()).To(BeNil())
	})

	Describe("Testing date/time properties", func() {
		var (
			timeNowOriginal func() time.Time
			now             time.Time
		)

		BeforeEach(func() {
			now = time.Now()

			timeNowOriginal = _timeNow
			_timeNow = func() time.Time {
				return now
			}
		})

		AfterEach(func() {
			_timeNow = timeNowOriginal
		})

		It("should return the current time", func() {
			Expect(fileInfo.ModTime()).To(Equal(now))
		})
	})
})
