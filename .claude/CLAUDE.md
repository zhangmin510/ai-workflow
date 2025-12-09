# 导入通用的AI Agent协作标准
@../AGENTS.md

# --- 以下是Claude Code专属的高级指令 ---

## Sub-agent定义
- 当需要进行安全审查时，请调用`security-reviewer` sub-agent。

## Hooks配置
- 在每次代码编辑后，自动运行`gofmt`。


# --- 团队共享规范 ---
- Go版本必须为 1.25+
- 所有错误处理必须使用 `github.com/pkg/errors`
...

# --- 个人偏好导入区 ---
# 每个团队成员可以将自己的个人配置文件放在这里
# 这个导入路径会被git忽略，但会被Claude Code加载
@~/.claude/personal-preferences.md
