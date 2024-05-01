[//]: # (Copyright [c] 2024 Pablo Sánchez)

[//]: # (Permission is hereby granted, free of charge, to any person obtaining a copy)
[//]: # (of this software and associated documentation files [the "Software"], to deal)
[//]: # (in the Software without restriction, including without limitation the rights)
[//]: # (to use, copy, modify, merge, publish, distribute, sublicense, and/or sell)
[//]: # (copies of the Software, and to permit persons to whom the Software is)
[//]: # (furnished to do so, subject to the following conditions:)

[//]: # (The above copyright notice and this permission notice shall be included in all)
[//]: # (copies or substantial portions of the Software.)

[//]: # (THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR)
[//]: # (IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,)
[//]: # (FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE)
[//]: # (AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER)
[//]: # (LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,)
[//]: # (OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE)
[//]: # (SOFTWARE.)

<img src="https://i.imgur.com/E4ocP1N.png" height="180em" style="border-left: 10px solid #ff9b22;"><img src="https://golang.org/doc/gopher/appenginegopher.jpg" height="180em">

# inbox-extract-data

App "inbox-extract-data" is a Go app to check a GMail inbox to analize Bolt e-scooter usages.

**TODO:**

  - Add more statistics
  - Get plans and travels over time
  - Add historical statistics
  - More test data

#  USE
## CREDENTIALS:

For inbox-extract-data to work you must have a gmail account and a file named "client_secret.json" containing your authorization info in the root directory of your project. To obtain credentials please see step one of this guide: https://developers.google.com/gmail/api/quickstart/go

 > Turning on the gmail API

 > - Use this wizard (https://console.developers.google.com/start/api?id=gmail) to create or select a project in the Google Developers Console and automatically turn on the API. Click Continue, then Go to credentials.
 
 > - On the Add credentials to your project page, click the Cancel button.
 
 > - At the top of the page, select the OAuth consent screen tab. Select an Email address, enter a Product name if not already set, and click the Save button.
 
 > - Select the Credentials tab, click the Create credentials button and select OAuth client ID.
 
 > - Select the application type Other, enter the name "Gmail API Quickstart", and click the Create button.
 
 > - Click OK to dismiss the resulting dialog.
 
 > - Click the file_download (Download JSON) button to the right of the client ID.
 
 > - Move this file to your working directory and rename it client_secret.json.

 > - Start the app and visit the link in the console (set the new app to production before)

 > - Copy the token from the URL and paste in console. Press enter.

## START:

```bash
go run main.go
```