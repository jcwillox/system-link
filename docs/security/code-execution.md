---
icon: material/code-braces
---

# Code Execution

* We don't allow arbitrary remote code execution by default. This means remote users can't just run any command on your system.
* Instead, we provide custom, buttons, switches, etc., with predefined actions, specified in your local configuration.
* This means you can also choose what commands to expose, for example, you might not want your server to be able to be remotely shutdown. You can choose to not expose that button entity, but could expose the lock button entity.
* This is a drastically more secure default, rather than allowing arbitrary code execution, which some other tools in this space do allow. Some other tools also don't allow you to disable basic power management commands, if you don't want or need remote users to be able to shut down, restart, sleep, etc., your device.
