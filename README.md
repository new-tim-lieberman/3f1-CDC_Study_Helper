# 3F1 CDC Study Helper

A simple study program for the 3F1 Career Development Course (CDC). It quizzes
you on Modules 1–4, tells you the right answer when you miss one, and
automatically brings missed questions back around in future sessions so they
stick.

No prior computer experience needed — just follow the steps for your
computer below.

## Getting the app

1. Go to the
   [Releases page](https://github.com/new-tim-lieberman/3f1-CDC_Study_Helper/releases)
   for this project.
2. Under the newest release, click **Assets** to expand the list of
   downloadable files.
3. Download the one file that matches your computer:
   - `3F1-CDC-Study-Windows.exe` — for Windows computers
   - `3F1-CDC-Study-Mac` — for Mac computers (works on both older Intel Macs
     and newer Apple Silicon Macs)

It's fine to move the downloaded file to your Desktop if that's easier to
find later. You do **not** need to download anything else from this page —
just that one file.

## Windows: how to run it

1. Find `3F1-CDC-Study-Windows.exe` wherever your browser saved it (usually
   your **Downloads** folder, or the Desktop if you moved it there), and
   **double-click it**.
2. Windows may show a blue box that says "Windows protected your PC." This is
   normal for a new program — click **More info**, then click **Run
   anyway**.
3. A black window will open with text in it. This is the app — leave it
   open and follow the instructions on screen.

## Mac: how to run it

Macs are stricter about running downloaded programs, so the first time takes
a few extra steps. After that, it's quick.

1. Open **Terminal**. Press `Command (⌘) + Space` to open Spotlight search,
   type `Terminal`, and press Enter. A window with text will open — this is
   normal, it's just a way to type commands to your computer.
2. Type `chmod +x ` (with a space after it), then **drag the downloaded
   `3F1-CDC-Study-Mac` file** (it's probably in your Downloads folder) into
   the Terminal window — this pastes its exact location. Press Enter. This
   step gives the file permission to run.
3. Type `xattr -d com.apple.quarantine ` (with a space), drag the same file
   into the Terminal window again, and press Enter. This tells your Mac it's
   okay to run a file you downloaded.
   If you see a message saying `No such xattr` — that's fine, it just means
   this step wasn't needed. Move on to the next step either way.
4. Now run the app: drag the same file into Terminal one more time (with
   nothing typed before it) and press Enter.
5. The app will start printing text in the same window — follow the
   instructions on screen.

Next time you want to study, you only need steps 1 and 4 (steps 2–3 are a
one-time thing).

## How to use it

Once it's running, you'll see a menu like this:

```
Available modules:
  1) Module 1 (79 questions)
  2) Module 2 (90 questions)
  3) Module 3 (40 questions)
  4) Module 4 (81 questions)
  0) All modules
  r) Review due questions (12 due)
  q) Quit
```

- Type a **number** (1, 2, 3...) and press Enter to study just that module.
- Type **0** to study every module at once.
- Type **r** to review questions you've missed before or haven't studied
  recently — this is the smartest way to study once you've done a session or
  two.
- Type **q** to quit.

For each question, type the **letter** (A, B, C...) or the **number** (1, 2,
3...) of your answer and press Enter. You'll be told right away if you got it
right, and shown the correct answer if you didn't.

At the end of a round, any questions you missed will be offered again so you
can retry them until you get them all correct.

### Your progress is saved automatically

Every time you answer a question, the app remembers whether you got it right.
Questions you keep getting right show up less often; questions you miss come
back sooner. You don't need to do anything — just close the app whenever
you're done (type `q` or close the window), and your progress will be there
next time you open it.

If you ever want to wipe your progress and start completely fresh, delete
this file:

- **Windows:** `C:\Users\<YourName>\.3f1_cdc_stats.json`
- **Mac:** `~/.3f1_cdc_stats.json` (your Home folder)

## For developers: running from source

If you have [Go](https://go.dev/dl/) installed and want to run the app
directly from source instead of the prebuilt binaries:

```
go run .
```

To build the binaries locally instead of relying on the release workflow
below:

```
GOOS=windows GOARCH=amd64 go build -o dist/3F1-CDC-Study-Windows.exe .
GOOS=darwin GOARCH=arm64 go build -o dist/3F1-CDC-Study-mac-arm64 .
GOOS=darwin GOARCH=amd64 go build -o dist/3F1-CDC-Study-mac-intel .
lipo -create -output dist/3F1-CDC-Study-Mac dist/3F1-CDC-Study-mac-arm64 dist/3F1-CDC-Study-mac-intel
```

`dist/` is gitignored — binaries aren't committed to the repo. Instead,
`.github/workflows/release.yml` builds both binaries and publishes them to
the [Releases page](https://github.com/new-tim-lieberman/3f1-CDC_Study_Helper/releases)
automatically whenever a version tag is pushed:

```
git tag v1.2.0
git push origin v1.2.0
```

That's the link the "Getting the app" section above sends non-technical
users to, so cut a tag (and get it merged/pushed) any time you want a new
build available for download.

To add more study questions, drop a new `moduleN.json` file into
`internal/quiz/data/` following the format of the existing files. You can
optionally add an `"explanation"` field to any question, which is shown when
that question is missed.
