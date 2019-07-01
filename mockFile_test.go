package mocka

import (
	"io"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("mockFile", func() {
	var file mockFile

	BeforeEach(func() {
		contents := "I am the file contents"
		file = mockFile{
			name:   "i_am_a_file",
			buf:    []byte(contents),
			offset: int64(0),
			base:   int64(0),
			limit:  int64(len(contents)),
		}
	})

	Describe("io.Closer", func() {
		It("implements oi.Closer", func() {
			mockFileType := reflect.TypeOf(&file)
			interfaceType := reflect.TypeOf((*io.Closer)(nil)).Elem()
			Expect(mockFileType.Implements(interfaceType)).To(BeTrue())
		})

		It("returns nil when the file is closed", func() {
			Expect(file.Close()).To(BeNil())
		})

	})

	Describe("io.Reader", func() {
		It("implements oi.Reader", func() {
			mockFileType := reflect.TypeOf(&file)
			interfaceType := reflect.TypeOf((*io.Reader)(nil)).Elem()
			Expect(mockFileType.Implements(interfaceType)).To(BeTrue())
		})

		It("reads the correct number of bytes", func() {
			bytesToRead := make([]byte, 4)

			_, err := file.Read(bytesToRead)

			Expect(bytesToRead).To(Equal([]byte("I am")))
			Expect(err).To(BeNil())

			_, err = file.Read(bytesToRead)

			Expect(bytesToRead).To(Equal([]byte(" the")))
			Expect(err).To(BeNil())
		})

		It("updates offest with number of bytes read", func() {
			bytesToRead := make([]byte, 2)

			n, err := file.Read(bytesToRead)

			Expect(err).To(BeNil())
			Expect(file.offset).To(Equal(int64(2)))
			Expect(n).To(Equal(2))

			n, err = file.Read(bytesToRead)

			Expect(err).To(BeNil())
			Expect(file.offset).To(Equal(int64(4)))
			Expect(n).To(Equal(2))
		})

		It("returns io.EOF error if you read all the bytes in the buffer", func() {
			bytesToRead := make([]byte, 2)
			file.offset = file.limit

			// attempt to read from the file again
			n, err := file.Read(bytesToRead)

			Expect(err).To(Equal(io.EOF))
			Expect(file.offset).To(Equal(int64(22)))
			Expect(n).To(Equal(0))
		})
	})

	Describe("io.Seeker", func() {
		It("implements oi.Seeker", func() {
			mockFileType := reflect.TypeOf(&file)
			interfaceType := reflect.TypeOf((*io.Seeker)(nil)).Elem()
			Expect(mockFileType.Implements(interfaceType)).To(BeTrue())
		})

		It("moves the pointer to a new byte position", func() {
			offest, err := file.Seek(10, io.SeekStart)

			Expect(offest).To(Equal(int64(10)))
			Expect(err).To(BeNil())

			offest, err = file.Seek(2, io.SeekStart)

			Expect(offest).To(Equal(int64(2)))
			Expect(err).To(BeNil())
		})

		It("moves the pointer to a new position relative to the current position", func() {
			offest, err := file.Seek(10, io.SeekCurrent)

			Expect(offest).To(Equal(int64(10)))
			Expect(err).To(BeNil())

			offest, err = file.Seek(2, io.SeekCurrent)

			Expect(offest).To(Equal(int64(12)))
			Expect(err).To(BeNil())
		})

		It("moves the pointer to a new position relative to the end of the buffer", func() {
			offest, err := file.Seek(-10, io.SeekEnd)

			Expect(offest).To(Equal(int64(12)))
			Expect(err).To(BeNil())

			offest, err = file.Seek(-2, io.SeekEnd)

			Expect(offest).To(Equal(int64(20)))
			Expect(err).To(BeNil())
		})

		It("returns error if the provided whence is invalid", func() {
			offest, err := file.Seek(10, 4)

			Expect(offest).To(Equal(int64(0)))
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("mocka: invalid whence received - 4"))
		})

		It("returns error if the offset is less than the file base", func() {
			offest, err := file.Seek(-1, io.SeekStart)

			Expect(offest).To(Equal(int64(0)))
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("mocka: desired offset -1 is below the offset base of 0"))
		})

		It("returns error if the offset is more than the file limit", func() {
			offest, err := file.Seek(1, io.SeekEnd)

			Expect(offest).To(Equal(int64(0)))
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("mocka: desired offset 23 is above the offset limit of 22"))
		})
	})

	Describe("io.ReaderAt", func() {
		It("implements oi.ReaderAt", func() {
			mockFileType := reflect.TypeOf(&file)
			interfaceType := reflect.TypeOf((*io.ReaderAt)(nil)).Elem()
			Expect(mockFileType.Implements(interfaceType)).To(BeTrue())
		})

		It("reads the correct number of bytes", func() {
			bytesToRead := make([]byte, 4)

			_, err := file.ReadAt(bytesToRead, file.offset)

			Expect(bytesToRead).To(Equal([]byte("I am")))
			Expect(err).To(BeNil())

			_, err = file.ReadAt(bytesToRead, file.offset+4)

			Expect(bytesToRead).To(Equal([]byte(" the")))
			Expect(err).To(BeNil())
		})

		It("returns the number of bytes read", func() {
			bytesToRead := make([]byte, 2)

			n, err := file.ReadAt(bytesToRead, file.offset)

			Expect(err).To(BeNil())
			Expect(n).To(Equal(2))
		})

		It("returns io.EOF error if you read all the bytes in the buffer", func() {
			bytesToRead := make([]byte, 1)

			n, err := file.ReadAt(bytesToRead, file.limit)

			Expect(err).To(Equal(io.EOF))
			Expect(n).To(Equal(0))
		})
	})

	Describe("io.Writer", func() {
		It("implements oi.Writer", func() {
			mockFileType := reflect.TypeOf(&file)
			interfaceType := reflect.TypeOf((*io.Writer)(nil)).Elem()
			Expect(mockFileType.Implements(interfaceType)).To(BeTrue())
		})

		It("writes bytes to file offset", func() {
			bytesWritten, err := file.Write([]byte("hahahaha"))

			Expect(err).To(BeNil())
			Expect(bytesWritten).To(Equal(8))
			Expect(file.buf).To(Equal([]byte("hahahaha file contents")))

			bytesWritten, err = file.Write([]byte("I am"))

			Expect(err).To(BeNil())
			Expect(bytesWritten).To(Equal(4))
			Expect(file.buf).To(Equal([]byte("I amhaha file contents")))
		})

		It("Returns io.ErrShortWrite error if write is longer than remaining bytes in the buffer", func() {
			bytesWritten, err := file.Write([]byte("This is the changed files contents that will not all make it"))

			Expect(err).ToNot(BeNil())
			Expect(err).To(Equal(io.ErrShortWrite))
			Expect(bytesWritten).To(Equal(22))
			Expect(file.buf).To(Equal([]byte("This is the changed fi")))

			bytesWritten, err = file.Write([]byte("I am"))
		})

	})

	Describe("os.File", func() {
		It("returns a nil for error", func() {
			_, err := file.Stat()

			Expect(err).To(BeNil())
		})

		It("returns fileInfo with the same filename", func() {
			fileInfo, _ := file.Stat()

			Expect(fileInfo).ToNot(BeNil())
			Expect(fileInfo.Name()).To(Equal(file.name))
		})

		It("returns fileInfo with same size", func() {
			fileInfo, _ := file.Stat()

			Expect(fileInfo).ToNot(BeNil())
			Expect(fileInfo.Size()).To(Equal(file.limit))
		})

		It("does not truncate the file if the size is larger than the file limit", func() {
			err := file.Truncate(40)

			Expect(err).To(BeNil())
			Expect(file.buf).To(Equal([]byte("I am the file contents")))
			Expect(file.limit).To(Equal(int64(22)))
		})

		It("does not truncate the file if the size is less than 0", func() {
			err := file.Truncate(-10)

			Expect(err).To(BeNil())
			Expect(file.buf).To(Equal([]byte("I am the file contents")))
			Expect(file.limit).To(Equal(int64(22)))
		})

		It("truncates the file to the size provided", func() {
			err := file.Truncate(8)

			Expect(err).To(BeNil())
			Expect(file.buf).To(Equal([]byte("I am the")))
			Expect(file.limit).To(Equal(int64(8)))
		})
	})

})
