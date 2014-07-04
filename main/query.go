package main

// Holds the query string to get public keys.

// No idea how much of this is necessary or correct.
// I take the view that Phab data might be insane, and re-validate it myself.
// This does NOT stop someone from throwing a 50 GB cat gif in their key body, or other creative ideas.
// I highly recommend doing your own homework before using this in a serious context.

// Some people, when confronted with a problem, think "I know, I'll use regular expressions..."


// From the MySQL documentation:
// Because MySQL uses the C escape syntax in strings (for example, “\n” to represent the newline character), you must double any “\” that you use in your REGEXP strings.
//
// Thus escaping a hyphen while in a regex character range is two slashes : \\-
// Buuut we're in golang, another system to pass escape chracters through (which are stacking like folding chairs), so QUAADRA-SLASH
//
const EscSeq = "\\\\"


// Regex to validate usernames, because why not.
//
// From an error message on a failed phabricator register attempt:
// Usernames must contain only numbers, letters, period, underscore and hyphen, and can not end with a period. They must have no more than 64 characters.
//
// I checked, and yes, `a` is a valid username. So truly 1-64 characters.
// I am 300% sure that this regex is 900% legit and not horrible at all.
const Regex_Username = "'^" +                                    // match start of line
			"[a-zA-Z0-9_" + EscSeq + "-" + EscSeq + ".]{1,63}" + // match 1 to 63 numbers, letters, underscores, hyphens, and periods.
			"[a-zA-Z0-9_" + EscSeq + "-]{0,1}" +                 // same thing, but 0 to 1 of the above excepting periods.
			"$'"                                                 // match end of line

// Silly, lazy, and wrong attempt to validate a chunk of base64
const Regex_KeyBody = "'^[a-zA-Z0-9" + EscSeq + "+/= ]+$'"


// STRING CATS ARE FAST CATS
const Query = "SELECT userName, isAdmin, keyBody, keyType, " + // the droids we're looking for
			"COALESCE(keyComment, '') " +                      // user comment, if any
			"FROM " + Table_SSH  + " AS ssh " +                // table with pubkeys in
			"JOIN " + Table_User + " AS u " +                  // table with users in
			"ON u.phid = ssh.userPHID " +                      // phabricator PHIDs are globally unique
			"WHERE " +
			"u.isDisabled = 0 AND " +                          // no disabled users!
			"u.username IS NOT NULL AND " +                    // lots of null checks
			"ssh.keybody IS NOT NULL AND " +                   // who knows if we need these
			"ssh.keyType IS NOT NULL AND " +                   // and who cares
			"u.username REGEXP " + Regex_Username + " AND " +  // validate usernames
			"ssh.keybody REGEXP " + Regex_KeyBody +            // validate usernames

			"ORDER BY ssh.id"                                  // consistent ordering by key ID, for change detection
