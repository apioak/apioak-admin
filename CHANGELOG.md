## 0.6.2
> Released on 2023.04.06

#### Feature

- Weaken resource association checks.
- Optimize association configuration synchronization when configuring configuration data operations.
- Optimizing routing upstream can be a null value operation.
- Cancel the operation restriction of the default route `/*`.



## 0.6.1
> Released on 2023.03.01

#### Feature

- Added upstream pool configuration function module.
- Added route self-selection and associated upstream function.


#### FIX

- Fixed the routing plugin publishing logic problem.
- Fixed the problem of remote data plane configuration synchronization exception.


#### Change

- Remove the node configuration in the routing configuration.
- Optimize the static document packaging system.
- Optimize remote data synchronization exception reminder.
- Optimize some parameter checks and logical association checks.
- The data table `oak_upstreams` adds fields `enable` and `release`.


#### Document

- Added change log documents, `CHANGELOG.md` and `CHANGELOG_CN.md`.
- Update documentation in `README.md`.
