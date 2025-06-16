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
