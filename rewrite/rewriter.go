// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package rewrite // import "WebParser/reader/rewrite"

import (
	"strconv"
	"strings"
	"text/scanner"

	"WebParser/logger"
	"WebParser/urllib"
)

type rule struct {
	name string
	args []string
}

// Rewriter modify item contents with a set of rewriting rules.
func Rewriter(entryURL string, content string, customRewriteRules string) string {
	rulesList := getPredefinedRewriteRules(entryURL)
	if customRewriteRules != "" {
		rulesList = customRewriteRules
	}

	rules := parseRules(rulesList)
	rules = append(rules, rule{name: "add_pdf_download_link"})

	logger.Debug(`[Rewrite] Applying rules %v for %q`, rules, entryURL)

	for _, rule := range rules {
		content = applyRule(entryURL, content, rule)
	}

	return content
}

func parseRules(rulesText string) (rules []rule) {
	scan := scanner.Scanner{Mode: scanner.ScanIdents | scanner.ScanStrings}
	scan.Init(strings.NewReader(rulesText))

	for {
		switch scan.Scan() {
		case scanner.Ident:
			rules = append(rules, rule{name: scan.TokenText()})

		case scanner.String:
			if l := len(rules) - 1; l >= 0 {
				text := scan.TokenText()
				text, _ = strconv.Unquote(text)

				rules[l].args = append(rules[l].args, text)
			}

		case scanner.EOF:
			return
		}
	}
}

func applyRule(entryURL string, content string, rule rule) string {
	switch rule.name {
	case "add_image_title":
		content = addImageTitle(entryURL, content)
	case "add_mailto_subject":
		content = addMailtoSubject(entryURL, content)
	case "add_dynamic_image":
		content = addDynamicImage(entryURL, content)
	case "add_youtube_video":
		content = addYoutubeVideo(entryURL, content)
	case "add_invidious_video":
		content = addInvidiousVideo(entryURL, content)
	case "add_youtube_video_using_invidious_player":
		content = addYoutubeVideoUsingInvidiousPlayer(entryURL, content)
	case "add_youtube_video_from_id":
		content = addYoutubeVideoFromId(content)
	case "add_pdf_download_link":
		content = addPDFLink(entryURL, content)
	case "nl2br":
		content = replaceLineFeeds(content)
	case "convert_text_link", "convert_text_links":
		content = replaceTextLinks(content)
	case "fix_medium_images":
		content = fixMediumImages(entryURL, content)
	case "use_noscript_figure_images":
		content = useNoScriptImages(entryURL, content)
	case "replace":
		// Format: replace("search-term"|"replace-term")
		if len(rule.args) >= 2 {
			content = replaceCustom(content, rule.args[0], rule.args[1])
		} else {
			logger.Debug("[Rewrite] Cannot find search and replace terms for replace rule %s", rule)
		}
	case "replace_title":
		// Format: replace_title("search-term"|"replace-term")
		// if len(rule.args) >= 2 {
		// 	Title = replaceCustom(Title, rule.args[0], rule.args[1])
		// } else {
		// 	logger.Debug("[Rewrite] Cannot find search and replace terms for replace rule %s", rule)
		// }
	case "remove":
		// Format: remove("#selector > .element, .another")
		if len(rule.args) >= 1 {
			content = removeCustom(content, rule.args[0])
		} else {
			logger.Debug("[Rewrite] Cannot find selector for remove rule %s", rule)
		}
	case "add_castopod_episode":
		content = addCastopodEpisode(entryURL, content)
	case "base64_decode":
		// if len(rule.args) >= 1 {
		// 	content = applyFuncOnTextcontent(content, rule.args[0], decodeBase64content)
		// } else {
		// 	content = applyFuncOnTextcontent(content, "body", decodeBase64content)
		// }
	case "parse_markdown":
		content = parseMarkdown(content)
	case "remove_tables":
		content = removeTables(content)
	case "remove_clickbait":
		// Title = removeClickbait(Title)
	}

	return content
}

func getPredefinedRewriteRules(entryURL string) string {
	urlDomain := urllib.Domain(entryURL)
	for domain, rules := range predefinedRules {
		if strings.Contains(urlDomain, domain) {
			return rules
		}
	}

	return ""
}
