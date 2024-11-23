package repository

type CompressRepository struct {
	compressedDir string
}

func NewCompressRepository(compressedDir string) *CompressRepository {
	return &CompressRepository{
		compressedDir: compressedDir,
	}
}

func (cr *CompressRepository) SaveFile(filename string) error {
	return nil
}

func (cr *CompressRepository) DeleteFile(filename string) error {
	return nil
}
