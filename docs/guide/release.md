# Release guidelines


odfe-cli versions are expressed as **x.y.z**, where **x** is the major version, **y** is the minor version, and **z** is the patch version, following [Semantic Versioning](https://semver.org/) terminology. 

**The major, minor & patch version is related to odfe-cli compatibility, instead of odfe’s compatibility.**


**Q: What goes into new patch version release train?**

* Known bugs that prevent users from executing tasks. 
    * eg: Failed to delete non-running detectors by delete command, failed to handle error.
* Security vulnerabilities. 
    * eg: leaking user credentials, insecure interaction with components.
* cli’s dependency is updated with new patch version which fixes bugs and security issues that affect your feature. (Treat your dependent library as your library)
    * eg: cobra released a new version which fixes some security vulnerabilities like [#1259](https://github.com/spf13/cobra/pull/1259)

Any committer, who fixed issues, **related to above use case**, should notify admin to initiate new patch release. There is no SLA for patch releases, since this is triggered based on severity of the issue.



**Q: What goes into new minor version release train?**

* Any new commands ( could be a new plugin, a new sub command for any plugin ). A new command could be a new API that was released long ago but being included in new odfe-cli release or new API that will be  released in next odfe release.
    * eg: on-board [k-nn](https://github.com/opendistro-for-elasticsearch/odfe-cli/pull/20/) plugin, add auto-complete feature for cli commands.
* New parameter/flags for any command.
    * eg: odfe-cli 1.0.0 only displays detector configuration using get command for given name, if user would also like to see [detector job](https://opendistro.github.io/for-elasticsearch-docs/docs/ad/api/#get-detector), they can add new flag (job) to get command to enable this feature.
* Any incompatible changes that were introduced with respect to API changes.
    * eg: if API is added in odfe 1.13.0 and updated in a backward incompatible way in later releases, this will be addressed in CLI as minor version

admin will trigger new minor release train whenever new plugin is on-boarded or every 45 days if any commits related to above use case is added. After every minor release, admin will create a new table where contributors can include candidates for next release.



**Q**: **What goes into new major version release train?**

* Any changes related to framework or addition of new component with respect to odfe-cli.
    * eg: Use new rest library (networking layer), framework re-design to ease new plugin on-boarding, user experience changes like update command formats, rename command names, etc. 

admin will trigger new major release train every 120 days if any commits related to above use case is added. After every major release, admin will create a new table related to next release where contributors can include candidates for next release.


