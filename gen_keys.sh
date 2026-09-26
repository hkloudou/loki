#!/usr/bin/env bash
# Regenerate keys.go from windows/lib/keys.h
set -euo pipefail
cd "$(dirname "$0")"
{
cat <<'HDR'
// Code generated from windows/lib/keys.h. DO NOT EDIT by hand.
// HID keyboard usage codes and modifier bitmasks (USB HID Usage Tables).
// Regenerate with: ./gen_keys.sh

package loki

// Modifier bitmask values for Keyboard.SetModifiers / Keyboard.Send.
const (
HDR
awk -F'[ \t]+' '/^#define KEY_MOD_/{
  name=$2; val=$3; gsub(/^KEY_MOD_/,"",name);
  n=split(tolower(name),parts,"_"); id="Mod";
  for(k=1;k<=n;k++){ p=parts[k]; id=id toupper(substr(p,1,1)) substr(p,2) }
  printf "\t%-10s = %s\n", id, val
}' windows/lib/keys.h
echo ")"
echo ""
echo "// HID keyboard usage codes for Keyboard.Type / Keyboard.Send."
echo "const ("
awk -F'[ \t]+' '/^#define KEY_/ && !/^#define KEY_MOD_/{
  name=$2; val=$3; gsub(/^KEY_/,"",name);
  n=split(tolower(name),parts,"_"); id="Key";
  for(k=1;k<=n;k++){ p=parts[k]; id=id toupper(substr(p,1,1)) substr(p,2) }
  printf "\t%-14s = %s\n", id, val
}' windows/lib/keys.h
echo ")"
} > keys.go
echo "regenerated keys.go"
