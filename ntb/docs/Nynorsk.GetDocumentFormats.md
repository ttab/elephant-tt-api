# Usage

`GetDocumentFormats()` returns the list of document formats supported by the translation service, including file suffixes, MIME types, and whether the content must be base64-encoded.

# Examples

Request:

```json
{}
```

Response (abbreviated):

```json
{
  "document_formats": [
    {
      "file_type": "html",
      "file_suffixes": [".html", ".htm"],
      "mime_types": ["text/html"],
      "base64_encode": false
    },
    {
      "file_type": "docx",
      "file_suffixes": [".docx"],
      "mime_types": ["application/vnd.openxmlformats-officedocument.wordprocessingml.document"],
      "base64_encode": true
    },
    {
      "file_type": "txt",
      "file_suffixes": [".txt"],
      "mime_types": ["text/plain"],
      "base64_encode": false
    },
    {
      "file_type": "json",
      "file_suffixes": [".json"],
      "mime_types": ["application/json"],
      "base64_encode": false
    }
  ]
}
```
