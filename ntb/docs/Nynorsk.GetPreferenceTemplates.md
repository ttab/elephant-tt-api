# Usage

`GetPreferenceTemplates()` returns the available preference templates and the language forms each template enables. Use a template name in the `prefs_template` field of a `Translate()` request instead of setting individual forms.

Available templates:

- **standard** — e-infinitiv and vi-pronoun only.
- **tradisjonell** — more traditional nynorsk forms.
- **moderne** — more modern, bokmål-adjacent forms.

# Examples

Request:

```json
{}
```

Response (abbreviated):

```json
{
  "preference_templates": {
    "moderne": {
      "values": [
        "tryggleik_sikkerheit.syn",
        "medviten_bevisst.syn",
        "infa_infe",
        "me_vi",
        "samd_einig.syn",
        "oppmode_oppfordre.syn"
      ]
    },
    "standard": {
      "values": [
        "infa_infe",
        "me_vi"
      ]
    },
    "tradisjonell": {
      "values": [
        "ansikt_andlet.syn",
        "enkelt_einskild.syn",
        "bli_verte.vb-bli2verte",
        "eigentleg_eigenleg.stav",
        "barna_borna.vok-o2a",
        "trost_trast.vok-o2a",
        "alter_altar.vok-a2e",
        "jern_jarn.vok-a2e",
        "kvalp_kvelp.vok-a2e",
        "spenn_span.vok-a2e"
      ]
    }
  }
}
```
