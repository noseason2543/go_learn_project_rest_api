package fileUsecases

import (
	"bytes"
	"context"
	"fmt"
	"go_learn_project_rest_api/config"
	"go_learn_project_rest_api/modules/files"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

type IFileUsecases interface {
	UploadToGCP([]*files.FileReq) ([]*files.FileRes, error)
	DeleteFileOnGCP([]*files.DeleteFileReq) error
}

type fileUsecases struct {
	cfg config.IConfig
}

func FileUsecases(cfg config.IConfig) IFileUsecases {
	return &fileUsecases{
		cfg: cfg,
	}
}

func (u *fileUsecases) UploadToGCP(req []*files.FileReq) ([]*files.FileRes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.newClient error: %v", err)
	}
	defer client.Close()

	jobsCh := make(chan *files.FileReq, len(req))
	resultCh := make(chan *files.FileRes, len(req))
	errCh := make(chan error, len(req))

	res := make([]*files.FileRes, 0)

	for _, r := range req {
		jobsCh <- r
	}
	close(jobsCh)

	numWorker := 5
	for i := 0; i < numWorker; i++ {
		go u.streamFileUpload(ctx, client, jobsCh, resultCh, errCh)
	}

	for a := 0; a < len(req); a++ {
		err := <-errCh
		if err != nil {
			return nil, err
		}

		result := <-resultCh
		res = append(res, result)
	}

	return res, nil
}

func (u *fileUsecases) streamFileUpload(ctx context.Context, client *storage.Client, jobs <-chan *files.FileReq, results chan<- *files.FileRes, errs chan<- error) {
	for job := range jobs {
		container, err := job.File.Open()
		if err != nil {
			errs <- err
			return
		}

		b, err := io.ReadAll(container)
		if err != nil {
			errs <- err
			return
		}

		buf := bytes.NewBuffer(b)
		// Upload an object with storage.Writer.
		wc := client.Bucket(u.cfg.App().GCPBucket()).Object(job.Destination + job.FileName).NewWriter(ctx)

		if _, err = io.Copy(wc, buf); err != nil {
			errs <- fmt.Errorf("io.Copy: %w", err)
			return
		}
		// Data can continue to be added to the file until the writer is closed.
		if err := wc.Close(); err != nil {
			errs <- fmt.Errorf("Writer.Close: %w", err)
			return
		}
		fmt.Printf("%v uploaded to %v.\n", job.FileName, job.Extension)

		newFile := &filesPub{
			file: &files.FileRes{
				FileName: job.FileName,
				Url:      fmt.Sprintf("https://storage.googleapis.com/%s/%s%s", u.cfg.App().GCPBucket(), job.Destination, job.FileName),
			},
			bucket:      u.cfg.App().GCPBucket(),
			destination: job.Destination,
		}

		if err := newFile.setPublicACL(ctx, client); err != nil {
			errs <- err
			return
		}

		errs <- nil
		results <- newFile.file
	}

}

type filesPub struct {
	bucket      string
	destination string
	file        *files.FileRes
}

func (f *filesPub) setPublicACL(ctx context.Context, client *storage.Client) error {
	acl := client.Bucket(f.bucket).Object(f.destination + f.file.FileName).ACL()
	if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return fmt.Errorf("object ACL set error: %w", err)

	}
	return nil
}

func (u *fileUsecases) DeleteFileOnGCP(req []*files.DeleteFileReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient : %v", err)
	}
	defer client.Close()

	jobsCh := make(chan *files.DeleteFileReq, len(req))
	errCh := make(chan error, len(req))

	for _, r := range req {
		jobsCh <- r
	}
	close(jobsCh)

	numWorker := 5
	for i := 0; i < numWorker; i++ {
		go u.deleteFile(ctx, client, jobsCh, errCh)
	}

	for a := 0; a < len(req); a++ {
		err := <-errCh
		if err != nil {
			return err
		}
	}

	return nil
}

// deleteFile removes specified object.
func (u *fileUsecases) deleteFile(ctx context.Context, client *storage.Client, jobs <-chan *files.DeleteFileReq, errs chan<- error) {
	for job := range jobs {
		o := client.Bucket(u.cfg.App().GCPBucket()).Object(job.Destination + job.FileName)

		// Optional: set a generation-match precondition to avoid potential race
		// conditions and data corruptions. The request to delete the file is aborted
		// if the object's generation number does not match your precondition.
		attrs, err := o.Attrs(ctx)
		if err != nil {
			errs <- fmt.Errorf("object.Attrs: %w", err)
			return
		}
		o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

		if err := o.Delete(ctx); err != nil {
			errs <- fmt.Errorf("Object(%q).Delete: %w", job.Destination+job.FileName, err)
		}
		fmt.Printf("Blob %v deleted.\n", job.Destination+job.FileName)
		errs <- nil
	}
}
