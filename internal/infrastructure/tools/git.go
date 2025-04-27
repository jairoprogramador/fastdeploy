package tools

import "unicode/utf8"
import "strings"

func ThereAreChanges() (bool, error) {
	message, err := ExecuteCommand("git", "status", "-s")
	if err != nil {
		return true, err
	}
	return utf8.RuneCountInString(message) > 0, nil
}

func GetCommitHash() (string, error) {
	commitHash, err := ExecuteCommand("git", "rev-parse", "HEAD")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(commitHash), nil
}

func GetCommitMessage(commitHash string) (string, error) {
	commitMessage, err := ExecuteCommand("git", "show", "-s", "--format=%s", commitHash)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(commitMessage), nil
}

func GetCommitAuthor(commitHash string) (string, error) {
	commitAuthor, err := ExecuteCommand("git", "show", "-s", "--format=%an <%ae>", commitHash)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(commitAuthor), nil
}