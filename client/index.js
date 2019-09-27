import AWSUtility from './aws_utility';

class Upload {
  awsUtility = null;
  constructor() {
    this.attachListener();
    this.awsUtility = new AWSUtility();
  }

  attachListener = () => {
    document.getElementsByTagName('form')[0].addEventListener('submit', (e) => {
      e.preventDefault();
      const fileupload = document.getElementById('fileupload');
      this.awsUtility.upload(fileupload.files[0], "image", () => {
        alert("uploaded");
      });
    });
  }
}

new Upload();