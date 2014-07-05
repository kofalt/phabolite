# Phabolite

A dead-simple way to port [Phabricator](http://phabricator.org) users' public keys into a [gitolite](http://gitolite.com/gitolite) repository.

Running a phabricator instance, but don't like their idea of git hosting layout?  
Run a gitolite instance, and take advantage of phab's web GUI for adding keys.

User-friendly and convenient!

## Setup

### Running the services

Setting up phabricator and gitolite is firmly out of scope; follow the links above!

### Adding database access

> Password hashes? In MY username table?  
> It may be more likely than you think.  
> Safely check your schema, FREE!

Phabolite needs access to the `user` table to get usernames & ID numbers, but that's also where password hashes go. Be smart; do not give phabolite carte blanche!

Below is an easy way to set up a user with the access it needs.

```sql
CREATE USER 'phabolite'@'localhost' IDENTIFIED BY 'terriblePassword';
GRANT SELECT ON phabricator_user.user_sshkey TO 'phabolite'@'localhost';
GRANT SELECT(userName, phid, isAdmin, isSystemAgent, isDisabled) ON phabricator_user.user TO 'phabolite'@'localhost';
FLUSH PRIVILEGES;
```

This should result in read-only access to select non-sensitive fields.

### Configure options

Phabolite desires a `phabolite.toml` file in the current directory.  
For help with MySQL connection strings, see [this documentation](https://github.com/go-sql-driver/mysql#dsn-data-source-name).

```toml
# MySQL connection string
server = "unix(/var/run/mysql/mysqld.sock)"

# MySQL account credentials
credentials = "phabolite:terriblePassword"

# SSH access string
ssh = "git@example.com"

# Loop forever?
loop = false

# How many seconds to wait between loops, if enabled
waitseconds = 30
```

## Future work

When dealing with a large number of users, the following issues my arise.  
I'd accept discussion tickets or pull requests regarding the following:

* Most notably, only one key per user is supported.
* The data validation bit was whipped up on the fly.
* No effort is made to detect if anything has changed; everything is flushed to disk, every time.
* Commit message is always "Phabolite generated update" rather than "Added one user", or similar.
