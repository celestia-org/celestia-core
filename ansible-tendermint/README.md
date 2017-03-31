# Ansible role for Tendermint

![Ansible plus Tendermint](img/a_plus_t.png)

* [Requirements](#requirements)
* [Variables](#variables)
* [Handlers](#handlers)
* [Example playbook that configures a Tendermint on Ubuntu](#example-playbook-that-configures-a-tendermint-on-ubuntu)

## Requirements

This role requires Ansible 2.0 or higher.

## Variables

Here is a list of all the default variables for this role, which are also
available in `defaults/main.yml`.

```
tendermint_version: 0.9.0
tendermint_archive: "tendermint_{{tendermint_version}}_linux_amd64.zip"
tendermint_download: "https://s3-us-west-2.amazonaws.com/tendermint/{{tendermint_version}}/{{tendermint_archive}}"
tendermint_download_folder: /tmp

tendermint_user: tendermint
tendermint_group: tendermint

# Upstart start/stop conditions can vary by distribution and environment
tendermint_upstart_start_on: start on runlevel [345]
tendermint_upstart_stop_on: stop on runlevel [!345]
tendermint_manage_service: true

tendermint_home: /opt/tendermint
tendermint_rpc_port: 46657
tendermint_proxy_app: tcp://127.0.0.1:46658

tendermint_log_file: /var/log/tendermint.log

tendermint_chain_id: mychain
tendermint_genesis_time: "{{ansible_date_time.iso8601_micro}}"
```

You can also change `templates/config.toml.j2` to suit your needs.

## Handlers

These are the handlers that are defined in `handlers/main.yml`.

* restart tendermint

## Example playbook that configures a Tendermint on Ubuntu

```
---

- hosts: all
  vars:
    tendermint_chain_id: MyAwesomeChain
    tendermint_seeds: "172.13.0.1:46656,172.13.0.2:46656,172.13.0.3:46656"
  roles:
    - ansible-tendermint
```

This playbook will install Tendermint and will create all the required
directories. But **it won't start the Tendermint if there are no validators in
genesis file**. See `templates/genesis.json.j2`.

You will need to collect validators public keys manually or using
`collect_public_keys.yml` given you have SSH access to all the nodes and add
them to `templates/genesis.json.j2`:

```
{
  "app_hash": "",
  "chain_id": "{{tendermint_chain_id}}",
  "genesis_time": "{{tendermint_genesis_time}}",
  "validators": [
    {
      "pub_key": [1, "3A4B5F5C34B19E5DBD2DC68E7D6FF7F46859A0657EDCA3274235A7EB127A0706"],
      "amount": 10,
      "name": "1"
    }
  ]
}
```

## Testing

```
vagrant up
```
