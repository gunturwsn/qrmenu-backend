package logging

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func HandlerInfo(c *fiber.Ctx, handler, msg string, statusCode int, resultCode string, kv ...interface{}) {
	fields := prependCodes(statusCode, resultCode, kv...)
	log.Printf("[handler] level=info request_id=%s handler=%s msg=%s%s",
		RequestIDFromFiber(c), handler, msg, formatFields(fields...))
}

func HandlerError(c *fiber.Ctx, handler, msg string, statusCode int, errCode string, err error, kv ...interface{}) {
	fields := prependCodes(statusCode, errCode, kv...)
	log.Printf("[handler] level=error request_id=%s handler=%s msg=%s%s error=%v",
		RequestIDFromFiber(c), handler, msg, formatFields(fields...), err)
}

func UsecaseInfo(scope, msg string, code string, kv ...interface{}) {
	fields := prependDomainCode(code, kv...)
	log.Printf("[usecase] level=info scope=%s msg=%s%s", scope, msg, formatFields(fields...))
}

func UsecaseError(scope, msg string, code string, err error, kv ...interface{}) {
	fields := prependDomainCode(code, kv...)
	log.Printf("[usecase] level=error scope=%s msg=%s%s error=%v", scope, msg, formatFields(fields...), err)
}

func RepoInfo(scope, msg string, code string, kv ...interface{}) {
	fields := prependDomainCode(code, kv...)
	log.Printf("[repository] level=info scope=%s msg=%s%s", scope, msg, formatFields(fields...))
}

func RepoError(scope, msg string, code string, err error, kv ...interface{}) {
	fields := prependDomainCode(code, kv...)
	log.Printf("[repository] level=error scope=%s msg=%s%s error=%v", scope, msg, formatFields(fields...), err)
}

func formatFields(kv ...interface{}) string {
	if len(kv) == 0 {
		return ""
	}
	max := len(kv) - len(kv)%2
	parts := make([]string, 0, (max/2)+1)
	for i := 0; i < max; i += 2 {
		parts = append(parts, fmt.Sprintf("%v=%v", kv[i], kv[i+1]))
	}
	if len(kv)%2 == 1 {
		parts = append(parts, fmt.Sprintf("arg%d=%v", len(kv)-1, kv[len(kv)-1]))
	}
	return " " + strings.Join(parts, " ")
}

func prependCodes(status int, code string, kv ...interface{}) []interface{} {
	fields := make([]interface{}, 0, len(kv)+4)
	if code != "" {
		fields = append(fields, labelForCode(status), code)
	}
	if status > 0 {
		fields = append(fields, "status_code", status)
	}
	fields = append(fields, kv...)
	return fields
}

func prependDomainCode(code string, kv ...interface{}) []interface{} {
	fields := make([]interface{}, 0, len(kv)+2)
	if code != "" {
		fields = append(fields, "code", code)
	}
	fields = append(fields, kv...)
	return fields
}

func labelForCode(status int) string {
	if status >= 400 {
		return "error_code"
	}
	return "result_code"
}
