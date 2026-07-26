package replynotification

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/models/reply"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func TestSetReplyUsesTheConvertedHTMLNotTheRawText(t *testing.T) {
	id := uuid.New()

	rep := &reply.Reply{
		ID:   id,
		Text: "I *love* <script>alert(1)</script> markdown",
		HTML: "<p>I <em>love</em> &lt;script&gt;alert(1)&lt;/script&gt; markdown</p>",
	}

	entity, err := New(
		&user.User{Username: "alice", Email: "alice@example.org", Language: "en"},
		"New reply",
	)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	entity.SetReply(rep, "bob", "_general/support/help#post-abc")

	if entity.Reply.ID != id {
		t.Errorf("ID = %v, want %v", entity.Reply.ID, id)
	}

	if got, want := entity.Reply.Text, rep.Text; got != want {
		t.Errorf("Text = %q, want %q", got, want)
	}

	if got, want := string(entity.Reply.HTML), rep.HTML; got != want {
		t.Errorf("HTML = %q, want the stored converted HTML %q", got, want)
	}

	if strings.Contains(string(entity.Reply.HTML), "<script>") {
		t.Error("HTML carries a live <script> tag; it must come from the " +
			"markdown-converted column, never from the raw reply text")
	}
}
