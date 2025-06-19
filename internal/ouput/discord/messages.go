package output

/*
General Messages
*/
const FAILURE_SERVER = "🤒 Something failed on our end, please submit a support ticket or reach out!"

const NO_ACTIVE_PROJECT = "No active project exists in this category, have an admin create one!"

const NOT_A_CHANNEL = "Hi! Sorry this bot can only be used in discord channels at the moment 😅"

const NOT_IN_A_CATEGORY = "Hello! This Bot can only be used within a category 😗"

const FAIL_INCORRECT_DATE = "❌ Incorrect date format, expect MM/DD/YYYY"

const FAIL_PERMS = "🔒 You lack the perms to use this command, Please reach out to a lead/admin!"

/*
Help comamnd messages
*/

const HELP_MSG = `>>> # Project Manager Bot
Thanks for using my MilestonePM bot 😊

For more in-depth guide on how to use the bot review its concepts here (holder):
[https://github.com/jaximus808/milePMBot]

**Docs Formatting**
*Arguments*: [argname]: <- command arguments, the [] indicate to put a value and argname is the name of the argument you are setting, where it becomes **argname:value**

*Example*: /help [cmd] <- The help command expects an arguemnt "cmd" with the command we indicate in []
valid usuage: /help cmd:project <- Indicates that the help should give more info about project commands

**Autocomplete**
- Please use autocomplete!!! This will give you the given options for a command and will speed up writing commands
- We'll continue working to add more autocomplete in the future 

*P.S.*: Spaces are okay with command arguments thanks to Discord's amzing API 😊
For these docs, we'll avoid spaces to reduce confusion 

## General Commands

` + "```/help ```" +
	`- prints this helper message

` + "```/help cmd:[command]```" + `
- This will give more info for a specific sub-command info
ex: /help cmd:project - gives information about project commands

## Project Commands

/project - handles projects settings and control (ADMIN ONLY)

## Milestone Commands

/milestone - handles the creation and movement of a milestone for a current project (ADMIN + LEADS ONLY)

## Task Commands

/task - handle creating, assigning, progress, completing, and approvals for tasks of a milestone 
`

// - ex: /project start msname:firstmilestone msdate:07/20/2026 msdesc:firstproject
const HELP_MSG_PROJECT = `>>> # Help Page: Project Commands

For more in-depth guide on how to use the bot review its concepts here (holder):
[https://github.com/jaximus808/milePMBot]

## Project Control (**ADMIN ONLY**)

` + "```/project start [msname] [msdate - in MM/DD/YYYY] [msdesc]```" + `
- Starts a project in the given discord channel category, and initiates the projects first milestone with the given values
- The initial milestone will become the active milestone

` + "```/project end [WIP]```" + `
- Ends a project and makes the active project nolonger apart of the active project
- Will ask for confirmation with a follow up message of the project's ref to prevent mistakes
- Will return a project ref you can use to get project info or resuming the project 

` + "```/project resume [ref] [WIP]```" + `
- Resumes a project that was ended and moves it into the current channels category
- Your discord account must have been an admin in that project

` + "```/project move [ref] [WIP]```" + `
- Moves an active project to another discord category
- All tasks, milestones, and roles will be maintained in this movement as well

## Project Settings (**ADMIN ONLY**)

` + "```/project set [setting] [setting]```" + `
- Updates the setting of a new project to have a specific value
- Our current settings include: changing output channel, update project desc, toggle sprints, sprint message, sprint length, toggle sprint pings

## Project Roles (**ADMIN ONLY**)

` + "```/project role [op] [user] [role - optional for remove]```" + `

## Project Info

` + "```/project info [WIP]```" + `
- Gives info for the current project 

`

const HELP_MSG_MILESTONE = `>>> # Help Page: Milestone Commands

For more in-depth guide on how to use the bot review its concepts here (holder):
[https://github.com/jaximus808/milePMBot]

## Milestone Control (**ADMIN**)

` + "```/milestone create [msname] [msdate] [msdesc]```" + `
- Creates a milestone for the project with the given milestone arguments
- Returns the milestone ref tied to this milestone
- Note: Does not make this milestone active

` + "```/milestone move [direction - next/prev]```" + `
- Moves the active milestone to be the previous or next milestone, which is determined by the due date 
- So if the current milestone is due 06/20/2025, and the next due milestone is 7/25/2025, then /milestone move direction:next will move the milestone the 7/25/2025 milestone

` + "```/milestone delete [taskref] [WIP]```" + `
- Removes the milestone and all tasks tied to it
- Will require a confirmation message after to confirm deletion
- This is permanent and can't be undone!

## Milestone Info

` + "```/milestone map```" + `
- Creates a milestone map that lists what tasks are unassigned, in progress, in review, and done for a milestone

`

const HELP_MSG_TASK = `>>> # Help Page: Task Commands

For more in-depth guide on how to use the bot review its concepts here (holder):
[https://github.com/jaximus808/milePMBot]

## Starting Tasks (**ADMIN + LEADS**)

` + "```/task create [name] [desc]```" + `
- Creates a new task for the current milestone with the task name and tas description
- Returns a task ref which you can use to refernece in the following commands

` + "```/task assign [user] [taskref] [expectation]```" + `
- Creates a new task for the current milestone with the task name and tas description
- Task expectations can both either use AGILE story points or due dates

## Task Progress

` + "```/task progress [taskref] [desc]```" + `
- Creates a progress update for the user's assigned task, with the description outlining the progress

` + "```/task compete [taskref] [desc]```" + `
- Marks a task as complete and ready for review, and will ping your assigner to review the task


## Reviewing Tasks (**Leads+Aadmins**)

` + "```/task approve [taskref]```" + `
- Marks a task as approved and done for the current milestone
- A Lead can only approve tasks they assigned to normal members
- Admins can approve any tasks, including tasks of other admins

` + "```/task reject [taskref] [desc]```" + `
- Marks a task as not approved and returned to in-progress 
- Desc is used to provide issues and feedback on what to fix
- A Lead can only reject tasks they assigned to normal members
- Admins can reject any tasks, including tasks of other admins

## General Task Commands

` + "```/task list```" + `
- Lists the tasks for the current milestone, and displays the status of all tasks

` + "```/task list [user]```" + `
- Lists the tasks that a user is assigned to, such as the tasks that are in progress, in review, and complete

`

/*
Project
*/

// PROJECT ERRORS
const FAIL_ALR_PROJECT = "❌ A project already exists here!"

// PROJECT SUCCESSES
const SUCCESS_CREATE_PROJECT = "# Your project has just been created! 🎉\n All commands now work within the channels in this category. \n Feel free to refer to /help or our webpage for references on my features 😃"

/*
Milestone failure messages
*/

const FAIL_MS_SAME_DATE = "❌ Two milestones can't have the same date"
const FAIL_ACTIVE_MS = "🤒 Failed to get the current milestone, please submit a support ticket or reach out!"

/*
Task failure
*/

// Failures
const FAIL_CREATE_TASK = "❌ Failed to create task, ensure this task has a unique name for the given milestone 😅"

const FAIL_TASK_DNE = "❌ Could not find that Task Ref, double check that task if for the current milestone 😅"

const ERROR_ARGS_ASSIGN = "❌ Expecting 3 args [@assign] [task_ref] [due_date or story points]!"

const ERROR_NOT_YOUR_TASK = "❌ This task isn't assigned to you 😶"

const ERROR_USER_NO_TASK = "❌ This user doesn't have any tasks assigned to them!"

// Success
const SUCCESS_APPROVING_TASK = "Yay! Task is now marked as approved :smile:"

const SUCCESS_COMPLETE_TASK = "🙌 Awesome! Marked as completed and sent to your assigner for review!"
const SUCCESS_TASK_LIST = "👀 Tasks List Successfully Made!"

const SUCCESS_PROGRESS_ADDED = "🫡 Got it! Updated progress and letting your assigner know"

const SUCCESS_REJECT = "👍 We'll mark this as not approved and notify the assigned person"
