
# Token Metadata Uploader (Go + Echo + Pinata)

This project provides a simple Go-based API that enables uploading files and metadata to IPFS via the Pinata API. It is designed to be used as part of a tokenization or NFT minting workflow.

## Technologies Used

- Go (Golang)
- Echo Web Framework
- Pinata API (for IPFS uploads)
- Docker


## Setup Instructions

### 1. Clone the repository

```bash
git clone https://github.com/your-username/go-deployment.git
cd go-deployment
```

### 2. Set up environment variables

Create a `.env` file or set the variables in your terminal session:

```
API_JWT=your_pinata_api_jwt
```

### 3. Run the application locally

```bash
go run main.go
```

The API will be available at: [http://localhost:8080](http://localhost:8080)

### 4. Test the `/upload` endpoint

You can test the endpoint using `curl`:

```bash
curl -X POST http://localhost:8080/upload \
  -F "name=Asset Name" \
  -F "description=Asset Description" \
  -F "json_name=Json Asset" \
  -F "file=@/path/to/image.png"
```

Sample response:

```json
{
  "tokenURI": "https://ipfs.io/ipfs/Qm..."
}
```

## Notes

- This service only handles metadata and file upload. Token minting is expected to be handled separately by a front-end or another backend service.
- Uploaded metadata follows the standard ERC-721 metadata structure.



