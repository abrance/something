package gz

import (
	"archive/tar"
	"fmt"
	"github.com/klauspost/compress/gzip"
	"io"
	"os"
	"strings"
)

// ReadFileInGzip
// @brief 直接读取tar.gz压缩包里面的某个目标文件而无需解压到本地
// @param srcFile 压缩包的路径
// @param FileInBall 文件在压缩包内的相对路径,非精确匹配,不能使用filepath.Join处理,否则无法适配windows系统
func ReadFileInGzip(srcFile, FileInBall string) ([]byte, error) {
	// 打开 tar.gz 文件
	file, err := os.Open(srcFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 创建一个 gzip.Reader
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()

	// 创建一个 tar.Reader
	tarReader := tar.NewReader(gzipReader)

	// 遍历 tar 文件
	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			// 到达文件末尾
			break
		}

		if err != nil {
			return nil, err
		}

		// 如果找到目标文件(因为打的包不规范,目录有./前缀,所以用contains来提高该函数的易用性)
		if strings.Contains(header.Name, FileInBall) {
			return io.ReadAll(tarReader)
		}
	}

	return nil, fmt.Errorf("未找到文件：%s\n", FileInBall)
}
