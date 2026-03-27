# Usage

`Translate()` translates a document between Scandinavian languages using Nynorskroboten. The document can be provided as a single string, a list of strings, or a map of named fragments.

The response returns the translated content in the same format as the request.

# Examples

Translate a simple text from bokmål to nynorsk:

```json
{
  "text": "Vi ønsker å rope hurra for pikens syvårsdag.",
  "file_type": "html",
  "source_language": "nb",
  "target_language": "nn"
}
```

```json
{
  "text": "Vi ønskjer å rope hurra for sjuårsdagen til jenta.",
  "guid": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

Translate multiple fragments as a list:

```json
{
  "texts": {
    "values": [
      "Regjeringen legger frem forslaget i dag.",
      "Opposisjonen er kritisk til endringene."
    ]
  },
  "file_type": "txt",
  "source_language": "nb",
  "target_language": "nn",
  "prefs_template": "moderne"
}
```

```json
{
  "texts": {
    "values": [
      "Regjeringa legg fram forslaget i dag.",
      "Opposisjonen er kritisk til endringane."
    ]
  },
  "guid": "bdf81353-c154-496f-b1c5-8ee1108809ca"
}
```

Translate named fragments as a map with individual language preferences:

```json
{
  "text_map": {
    "entries": {
      "headline": "Statsministeren holder tale",
      "body": "Regjeringen la i dag frem sitt forslag til ny lov."
    }
  },
  "file_type": "txt",
  "source_language": "nb",
  "target_language": "nn",
  "prefs_template": "tradisjonell",
  "prefs": {
    "samtidig_samstundes": {
      "enabled": true
    }
  }
}
```

```json
{
  "text_map": {
    "entries": {
      "body": "Regjeringa la i dag fram framlegget sitt til ny lov.",
      "headline": "Statsministeren held tale"
    }
  },
  "guid": "504ebf04-5f0e-46da-8e99-dd97bf86b17f"
}
```
