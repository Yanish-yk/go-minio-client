package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gofilegttptest/common"
	"gofilegttptest/fileio"
	"io/ioutil"
	"log"
	"strconv"
	"time"
)
func main() {
	ctx := context.Background()
	configFileLocation := flag.String("configfile", "configuration.toml", "configuration file")
	config,err := common.LoadTomlConfig(*configFileLocation)
	if err != nil  {
		log.Println(err)
		return
	}
	dir_name:=fileio.Getfilename(config.Bucket.FilePath)
	var  last int
	if len(dir_name)>3 {
		last=2
	} else {
		last=len(dir_name)
	}
	//var bucket_name  []string
	dir_name=dir_name[len(dir_name)-last:len(dir_name)]
	fmt.Println(dir_name)
	//桶名称加公司名称
	//temp:=common.Filenameaddcompany(dir_name,config.Bucket.Company)
	fmt.Println(dir_name,"dir_name")
	//初始化客户端
	minioClient, err := minio.New(config.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.Minio.AccessKeyID, config.Minio.SecretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}
	for i, _ := range dir_name {
		found, err := minioClient.BucketExists(ctx,config.Bucket.Company)
		if err != nil {
			fmt.Println(err)
			return
		}
		if !found {
			err = minioClient.MakeBucket(ctx,config.Bucket.Company,minio.MakeBucketOptions{Region: ""})
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Successfully created mybucket.")
			GetAllFile(config.Bucket.FilePath+"/"+dir_name[i],*minioClient,ctx,config.Bucket.Company,dir_name[i])
		}else {
			GetAllFile(config.Bucket.FilePath+"/"+dir_name[i],*minioClient,ctx,config.Bucket.Company,dir_name[i])
		}
	}
	//lists, err := minioClient.ListBuckets(ctx)
	//
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//for _, list := range lists {
	//	fmt.Println(list.Name)
	//}
	//下载对象文件
	//err = minioClient.FGetObject(ctx, "ibafile", "dockerimages", "/data/qwe",minio.GetObjectOptions{})
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println("下载")
}

func GetAllFile(pathname string,minioClient minio.Client,ctx context.Context,bucketname string,objectname string) error {
	year:=time.Now().Year()
	stringyear := strconv.Itoa(year)

	rd, err := ioutil.ReadDir(pathname)
	for _, fi := range rd {
		if fi.IsDir() {
			fmt.Printf("[%s]\n", pathname+"\\"+fi.Name())
			GetAllFile(pathname,minioClient,ctx,bucketname,objectname)
		} else {
			_, err := minioClient.FPutObject(ctx, bucketname,stringyear+"/"+objectname+"/"+fi.Name(),pathname+"/"+fi.Name(), minio.PutObjectOptions{ContentType: ""})
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println("----------以上传：",fi.Name(),"-----------------")
		}
	}
	return err
}