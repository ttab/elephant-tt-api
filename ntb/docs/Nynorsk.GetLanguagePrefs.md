# Usage

`GetLanguagePrefs()` returns the available optional language forms. Each form can be enabled or disabled in the `prefs` field of a `Translate()` request. Some forms support word-level overrides via the `modify_for` list.

# Examples

Request:

```json
{}
```

Response (abbreviated):

```json
{
  "language_prefs": {
    "aftan_eftan.vok-a2e": {
      "category": "Lydendringar",
      "description": "aftan → eftan"
    },
    "allmenn_ålmenn.vok-a2å": {
      "category": "Lydendringar",
      "description": "allmenn → ålmenn, klar → klår, skap → skåp, trav → tråv ofl."
    },
    "alter_altar.vok-a2e": {
      "category": "Lydendringar",
      "description": "alter → altar"
    },
    "ankel_okle.syn": {
      "category": "Synonym og stavemåtar",
      "description": "ankel → okle"
    }
  }
}
```
