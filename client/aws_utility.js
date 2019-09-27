export default class AWSUtility {
  PART_SIZE = 100 * 1024 * 1024; // 100MB // Minimum part size defined by aws s3 is 5 MB, maximum 5 GB
  file;
  fileInfo = {};
  sendBackData = null;

  // file fileObject
  upload = (file, bucketType, callback) => {
    this.file = file;
    this.fileInfo = {
      name: this.file.name,
      type: this.file.type,
      size: this.file.size,
    };
    this.sendBackData = null;
    this.createMultipartUpload(bucketType, response => {
      this.sendBackData = response.result;
      this.uploadMultipartFile(bucketType, callback);
    });
  }

  createMultipartUpload = (bucketType, callback) => {
    const reqParams = {
      ContentType: this.fileInfo.type,
      BucketType: bucketType,
    };
    this.postData(
      'AWSService.StartUpload',
      reqParams,
      err => console.dir(err),
      res => callback(res)
    );
  }

  uploadMultipartFile = async (bucketType, callback) => {
    try {
      const NUM_CHUNKS = Math.floor(this.fileInfo.size / this.PART_SIZE) + 1
      let promisesArray = []

      for (let index = 1; index < NUM_CHUNKS + 1; index++) {
        let start, end, blob;
        start = (index - 1) * this.PART_SIZE
        end = (index) * this.PART_SIZE
        blob = (index < NUM_CHUNKS) ? this.file.slice(start, end) : this.file.slice(start)

        const urlResponse = await this.getUploadUrl({
          MediaID: this.sendBackData.MediaID,
          PartNumber: index,
          UploadId: this.sendBackData.UploadID,
          BucketType: bucketType,
        });
        const promise = this.uploadPart(urlResponse.result.URL, blob);
        promisesArray.push(promise);
      }

      let resolvedArray = await Promise.all(promisesArray)
      console.log(resolvedArray, ' resolvedAr')

      let uploadPartsArray = []
      resolvedArray.forEach((resolvedPromise, index) => {
        uploadPartsArray.push({
          ETag: resolvedPromise.headers.get('etag'),
          PartNumber: index + 1
        })
        console.log(resolvedPromise);
        console.log(index);
      })

      this.completeUpload(uploadPartsArray, bucketType,  callback);
    } catch (err) {
      console.log(err)
    }
  }

  getUploadUrl = (params) => {
    return new Promise((resolve, reject) => {
      this.postData(
        'AWSService.GetUploadURL',
        params,
        err => reject(err),
        res => resolve(res)
      );
    });
  }

  uploadPart = (url, blob) => {
    return this.putData(url, blob, this.fileInfo.type)
  }

  completeUpload = (uploadPartsArray, bucketType, callback) => {
    const params = {
      MediaID: this.sendBackData.MediaID,
      Parts: uploadPartsArray,
      UploadId: this.sendBackData.UploadID,
      BucketType: bucketType,
    }
    this.postData(
      'AWSService.CompleteUpload',
      params,
      err => console.dir(err),
      res => callback(res)
    );
  }

  postData = (methodName, data, onError, onSuccess) => {
    const apiRequestData = {
      method: methodName,
      id: 1,
      params: [data],
    };
    const responsePromise = fetch(`http://localhost:10021/api/`, {
      method: "POST",
      headers: {
        'Content-Type': 'application/json',
      },
      // credentials: 'same-origin',
      body: JSON.stringify(apiRequestData),
    });
    this.onResponseReceipt(responsePromise, onError, onSuccess);
  }

  onResponseReceipt = (responsePromise,  onError, onSuccess) => {
    responsePromise.then(response => response.json())
      .then(responseJson => {
        if (responseJson && onSuccess) {
          onSuccess({ result: responseJson.result, error: null });
        } else if (onError) {
          if (!responseJson) {
            onError({ result: null, error: 'No response received.' });
          } else if (!responseJson.error === false) {
            onError({ result: null, error: responseJson.error });
          } else if (!responseJson.result) {
            onError({ result: null, error: 'No response received.' });
          }
        }
      })
      .catch(exception => {
        
        onError({ result: null, error: exception });
      });
  }

  putData = (url, data) => {
    return fetch(url, {
      method: "PUT",
      body: data,
    })
  }
}