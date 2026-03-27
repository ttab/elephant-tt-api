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

Translate an HTML article from bokmål to nynorsk:

```json
{
  "guid": "519fd7b7-8bc3-4c35-81c8-6f037ce92c35",
  "text": "<h1>Amerikanske sedler får Trumps signatur</h1><p>Donald Trump bryter med en 165 år gammel tradisjon når han blir den første sittende presidenten som får signaturen sin på amerikanske pengesedler.</p><p>Den første seddelen med Trumps navnetrekk blir 100-dollarseddelen, som oppdateres til sommeren, i forbindelse med USAs 250-årsjubileum.</p><p>Tidligere er det navnetrekkene til finansministeren og tjenestemannen med ansvar for USAs gullreserver som har prydet amerikanske pengesedler. Sistnevnte blir nå fjernet for første gang siden 1861 for å gi plass til Trumps signatur.</p>",
  "file_type": "html",
  "source_language": "nb",
  "target_language": "nn"
}
```

```json
{
  "text": "<h1>Amerikanske setlar får Trumps signatur</h1><p>Donald Trump bryt med ein 165 år gammal tradisjon når han blir den første sitjande presidenten som får signaturen sin på amerikanske pengesetlar.</p><p>Den første setelen med Trumps namnetrekk blir 100-dollarsetelen, som blir oppdatert til sommaren, i samband med USAs 250-årsjubileum.</p><p>Tidlegare er det namnetrekka til finansministeren og tenestemannen med ansvar for USAs gullreservar som har pryda amerikanske pengesetlar. Sistnemnde blir no fjerna for første gong sidan 1861 for å gi plass til Trumps signatur.</p>",
  "guid": "519fd7b7-8bc3-4c35-81c8-6f037ce92c35"
}
