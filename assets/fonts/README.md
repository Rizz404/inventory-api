# Font Setup for PDF Export (Japanese Support)

## Required Fonts

For Japanese language support in PDF exports, you need to download **Noto Sans JP** fonts:

### Download Instructions

1. Visit: https://fonts.google.com/noto/specimen/Noto+Sans+JP
2. Click "Download family" button
3. Extract the downloaded ZIP file
4. Copy the following files to this directory (`assets/fonts/`):
   - `NotoSansJP-Regular.ttf`
   - `NotoSansJP-Bold.ttf`

### Alternative: Direct Download

```bash
# Using curl (Linux/Mac)
cd assets/fonts
curl -L "https://github.com/notofonts/noto-cjk/raw/main/Sans/OTF/Japanese/NotoSansJP-Regular.otf" -o NotoSansJP-Regular.ttf
curl -L "https://github.com/notofonts/noto-cjk/raw/main/Sans/OTF/Japanese/NotoSansJP-Bold.otf" -o NotoSansJP-Bold.ttf
```

### PowerShell (Windows)

```powershell
cd assets\fonts
Invoke-WebRequest -Uri "https://github.com/notofonts/noto-cjk/raw/main/Sans/OTF/Japanese/NotoSansJP-Regular.otf" -OutFile "NotoSansJP-Regular.ttf"
Invoke-WebRequest -Uri "https://github.com/notofonts/noto-cjk/raw/main/Sans/OTF/Japanese/NotoSansJP-Bold.otf" -OutFile "NotoSansJP-Bold.ttf"
```

## Supported Languages

- ✅ English (default fonts)
- ✅ Japanese (requires Noto Sans JP - 日本語対応)
- 🔜 Indonesian (default fonts work)

## File Structure

```
assets/
└── fonts/
    ├── README.md              # This file
    ├── NotoSansJP-Regular.ttf # Download this
    └── NotoSansJP-Bold.ttf    # Download this
```

## Notes

- If Japanese fonts are not found, the system will fallback to default fonts
- Default fonts do not support CJK (Chinese, Japanese, Korean) characters
- CJK characters will display as "..." without proper fonts
- Font files are ~2-4MB each

## License

Noto Sans JP is licensed under the SIL Open Font License (OFL)
https://scripts.sil.org/OFL
