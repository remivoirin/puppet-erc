# Puppet ERC
RESTful API written in Go to manages roles in a Puppet role-profile pattern. ERC stands for "External Role Classifier".

## Why an External Role Classifier?

Puppet has typically two sources of data: variables set in the Puppet code, and Hiera. Many organizations use a role-profile pattern to help managing node hierarchy and code inheritance, and there are lots of ways to implement it.

A common problem is to define how Puppet agents get their role name, while the role content (included profiles or modules) is provided by the code or Hiera.

This ERC allows you to add a "role" fact to your Puppet agents, so you can define roles contents using the "$role" variable (either inside a manifest, module or an additional Hiera datadir).

The classifier works with entries in a DB, which you can manage with CRUD methods. An entry is made of three parts:
* The hostname regexp, which represent all Puppet agent hostnames that will match the role.
* The role.
* An additional comment for documentation purposes.

## Requirements

* A fully working Puppet setup.
* Pluginsync must be enabled on the Puppet server (Puppet 3.4+).
* curl on Puppet agents. 

## Deployment

### Server

**Clone the repo**
```sh
git clone https://github.com/remivoirin/puppet-erc.git && cd puppet-erc
```

**Launch**
* Using Go (requires the Echo framework)
```sh
cd code
go run *.go
```
You might want to use systemd or supervisord to manage the process execution as the access log is written on the standard output.

### Agents

* Customize the SERVER variable in facts.d/role.sh.
* Copy the fact to the facts.d/ directory of your base module (or any module that is included on every agent).

## Usage

**Insert an entry**: use any Go regexp for host_regex.
```sh
curl -X PUT -H 'Content-Type: application/json' -d '{ "host_regex": "test([0-9]+).domain.tld", "role": "testrole", "comment": "Test comment" }' http://yourhost:14002/insert
```

**List all roles**
```sh
curl http://yourhost:14002/list | jq .
```

**Get role for a host**: iterates on every hostname regexp, returns the role of the first matching regex, default otherwise. Iteration is made on a "latest inserted, first checked" basis.
* Fulltext (used by Puppet fact)
```sh
curl http://yourhost:14002/role/fulltext/myhostname.mydomain.tld
```

**Delete**
* with an id
```sh
curl -X DELETE http://yourhost:14002/delete/id/15
```
Other methods for deletion will be available. For now, list entries and delete the corresponding id.

## Known bugs and things to fix
* TODO: Provide a binary/static version (Github release).
* TODO: Make the "get hosts for a role" functions.
* TODO: Delete every host regex matching a role (= delete a role).
* TODO: Delete an entry matching an host regex.
* TODO: Check for any duplicate entry while inserting a new one.
* TODO: Add some auth on insert and delete.
* TODO: Rewrite the role fact in Ruby.

## Contact
I am @lagsfr on the Puppet community channel, @LagWire on Twitter, and my email is remi _at_ lags.is.

## License
MIT. See the LICENSE file for details.

## References and credits
* Using echo framework https://github.com/labstack/echo (MIT License)
* Using go-sqlite3 https://github.com/mattn/go-sqlite3 (MIT License)
