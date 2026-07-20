package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

var (
	flagBoard    string
	flagSet      string
	flagOut      string
	flagChrome   string
	flagPassword string
	flagQuality  int
	flagWidth    int
	flagHeight   int
	flagOnly     string
	flagClip     string
	flagResize   string
)

func init() {
	flag.StringVar(&flagBoard, "board", "http://127.0.0.1:3100", "base URL of the board to shoot")
	flag.StringVar(&flagSet, "set", "", "shoot only this set of shots, empty for every set")
	flag.StringVar(&flagOut, "out", "", "directory the screenshots are written into, overriding the set's own")
	flag.StringVar(&flagChrome, "chrome", "", "path to a Chrome or Chromium binary")
	flag.StringVar(&flagPassword, "password", "hyperhyper!", "password of the seeded users")
	flag.IntVar(&flagQuality, "quality", 0, "webp quality, 0 to 100, overriding the set's own")
	flag.IntVar(&flagWidth, "width", 0, "viewport width, overriding the set's own")
	flag.IntVar(&flagHeight, "height", 0, "viewport height, overriding the set's own")
	flag.StringVar(&flagOnly, "only", "", "shoot only the shots whose file path contains this substring")
	flag.StringVar(&flagClip, "clip", "", "element to clip every shot to, empty for the viewport, overriding the set's own")
	flag.StringVar(&flagResize, "resize", "", "size every shot is resized to, empty for none, overriding the set's own")
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "screenshots: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	chrome, err := findChrome()
	if err != nil {
		return err
	}

	if _, err = exec.LookPath("magick"); err != nil {
		return fmt.Errorf("ImageMagick is required to write webp: %w", err)
	}

	chosen, err := chosenSets()
	if err != nil {
		return err
	}

	shots := 0

	for _, s := range chosen {
		n, err := shoot(chrome, s)
		shots += n

		if err != nil {
			return fmt.Errorf("%s: %w", s.Name, err)
		}
	}

	fmt.Printf("==> wrote %d screenshots\n", shots)

	return nil
}

func chosenSets() ([]set, error) {
	given := map[string]bool{}
	flag.Visit(func(f *flag.Flag) { given[f.Name] = true })

	var chosen []set

	for _, s := range sets {
		if flagSet != "" && s.Name != flagSet {
			continue
		}

		chosen = append(chosen, s.override(given))
	}

	if len(chosen) == 0 {
		return nil, fmt.Errorf("there is no set named %q", flagSet)
	}

	return chosen, nil
}

func (s set) override(given map[string]bool) set {
	if given["out"] {
		s.Out = flagOut
	}
	if given["width"] {
		s.Width = flagWidth
	}
	if given["height"] {
		s.Height = flagHeight
	}
	if given["clip"] {
		s.Clip = flagClip
	}
	if given["resize"] {
		s.Resize = flagResize
	}
	if given["quality"] {
		s.Quality = flagQuality
	}

	return s
}

func shoot(chrome string, s set) (int, error) {
	var todo []shot

	for _, sh := range s.Shots {
		if flagOnly == "" || strings.Contains(sh.File, flagOnly) {
			todo = append(todo, sh)
		}
	}

	if len(todo) == 0 {
		return 0, nil
	}

	fmt.Printf("==> %s, into %s\n", s.Name, s.Out)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(chrome),
		chromedp.WindowSize(s.Width, s.Height),
		chromedp.Flag("hide-scrollbars", true),
		chromedp.Flag("force-color-profile", "srgb"),
		chromedp.Flag("font-render-hinting", "none"),
	)

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancelAlloc()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancelTimeout := context.WithTimeout(ctx, 10*time.Minute)
	defer cancelTimeout()

	if err := chromedp.Run(ctx); err != nil {
		return 0, fmt.Errorf("starting %s: %w", chrome, err)
	}

	if s.Clip == "" {
		if err := chromedp.Run(ctx,
			chromedp.EmulateViewport(int64(s.Width), int64(s.Height)),
		); err != nil {
			return 0, fmt.Errorf("sizing the viewport: %w", err)
		}
	}

	current := "\x00"
	shots := 0

	for _, sh := range todo {
		if sh.As != current {
			if err := signIn(ctx, sh.As); err != nil {
				return shots, fmt.Errorf("signing in as %q: %w", sh.As, err)
			}
			current = sh.As
		}

		if err := capture(ctx, s, sh); err != nil {
			return shots, fmt.Errorf("%s: %w", sh.File, err)
		}

		shots++
		fmt.Printf("    %s\n", sh.File)
	}

	return shots, nil
}

func signIn(ctx context.Context, username string) error {
	if err := chromedp.Run(ctx,
		chromedp.Navigate(flagBoard+"/session/signout"),
		chromedp.WaitReady("body", chromedp.ByQuery),
	); err != nil {
		return err
	}

	if username == "" {
		return nil
	}

	return chromedp.Run(ctx,
		chromedp.Navigate(flagBoard+"/session/signin"),
		chromedp.WaitVisible(`input[name="username"]`, chromedp.ByQuery),
		chromedp.SendKeys(`input[name="username"]`, username, chromedp.ByQuery),
		chromedp.SendKeys(`input[name="password"]`, flagPassword, chromedp.ByQuery),
		chromedp.Click(`button[type="submit"]`, chromedp.ByQuery),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.Sleep(200*time.Millisecond),
	)
}

func capture(ctx context.Context, s set, sh shot) error {
	actions := []chromedp.Action{
		chromedp.Navigate(flagBoard + sh.URL),
		chromedp.WaitReady("body", chromedp.ByQuery),
	}

	if sh.Wait != "" {
		actions = append(actions,
			chromedp.WaitVisible(sh.Wait, chromedp.ByQuery))
	}

	clip := s.Clip
	if sh.Clip != "" {
		clip = sh.Clip
	}

	below := s.Below
	if sh.Below != "" {
		below = sh.Below
	}

	var top float64
	if clip == "" && below != "" {
		actions = append(actions, fitBelow(s, below, &top))
	}

	actions = append(actions, sh.Prep...)
	actions = append(actions, chromedp.Sleep(250*time.Millisecond))

	var buf []byte
	switch {
	case clip != "":
		actions = append(actions,
			chromedp.Screenshot(clip, &buf, chromedp.ByQuery))
	case below != "":
		actions = append(actions, captureBelow(s, &top, &buf))
	default:
		actions = append(actions, chromedp.CaptureScreenshot(&buf))
	}

	if err := chromedp.Run(ctx, actions...); err != nil {
		return err
	}

	return writeWebP(filepath.Join(s.Out, sh.File), buf, s)
}

func fitBelow(s set, sel string, top *float64) chromedp.Action {
	bottom := fmt.Sprintf(`(() => {
		const e = document.querySelector(%q);
		if (!e) return -1;
		return e.getBoundingClientRect().bottom + window.scrollY;
	})()`, sel)

	return chromedp.ActionFunc(func(ctx context.Context) error {
		if err := chromedp.Evaluate(bottom, top).Do(ctx); err != nil {
			return err
		}

		if *top < 0 {
			return fmt.Errorf("there is no %q to shoot below", sel)
		}

		return chromedp.EmulateViewport(
			int64(s.Width), int64(*top)+int64(s.Height)).Do(ctx)
	})
}

func captureBelow(s set, top *float64, buf *[]byte) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		shot, err := page.CaptureScreenshot().
			WithFormat(page.CaptureScreenshotFormatPng).
			WithCaptureBeyondViewport(true).
			WithFromSurface(true).
			WithClip(&page.Viewport{
				X:      0,
				Y:      *top,
				Width:  float64(s.Width),
				Height: float64(s.Height),
				Scale:  1,
			}).Do(ctx)
		if err != nil {
			return err
		}

		*buf = shot

		return nil
	})
}

func writeWebP(dst string, png []byte, s set) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	args := []string{"png:-", "-strip"}

	if s.Resize != "" {
		args = append(args, "-resize", s.Resize)
	}

	args = append(args,
		"-define", "webp:method=6",
		"-define", "webp:thread-level=1",
		"-quality", strconv.Itoa(s.Quality),
		"webp:"+dst,
	)

	cmd := exec.Command("magick", args...)
	cmd.Stdin = bytes.NewReader(png)
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func findChrome() (string, error) {
	if flagChrome != "" {
		return flagChrome, nil
	}

	if env := os.Getenv("CHROME"); env != "" {
		return env, nil
	}

	for _, name := range []string{
		"chromium",
		"chromium-browser",
		"google-chrome",
		"google-chrome-stable",
	} {
		if path, err := exec.LookPath(name); err == nil {
			return path, nil
		}
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	matches, _ := filepath.Glob(filepath.Join(home,
		".cache", "ms-playwright", "chromium-*", "chrome-linux64", "chrome"))
	if len(matches) > 0 {
		return matches[len(matches)-1], nil
	}

	return "", fmt.Errorf(
		"no Chrome or Chromium found, pass -chrome or set CHROME")
}
