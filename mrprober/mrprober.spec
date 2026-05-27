Name:           mrprober
Version:        1.0
Release:        1%{?dist}
Summary:        Probe scheduler and reporter for K8S.

License:        MIT License
URL:            https://mr-prober.internal/
Source0:        %{name}-%{version}.tar.gz

BuildRequires:  golang >= 1.21
BuildRequires:  git
BuildRequires:  make

ExclusiveArch:  amd64 x86_64

%description
DaemonSet Utility for reporting various host probes.


%global debug_package %{nil}


%prep
%autosetup


%build
make


%install
install -Dpm 0755 %{name} %{buildroot}%{_bindir}/%{name}


%files
%license        LICENSE
%doc            README.md
%{_bindir}/%{name}


%changelog
* Thu Nov 02 2023 nikita
- bump version to 1.0
- bug fixes
- include dashboard to k8s deployment
- add meta probe as a metric
* Tue Oct 31 2023 nikita
- bump version to 0.3
- fix duplicated metrics with unordered labels
- updated metric names and switched to labels
- rename succ -> pass tag in raw output
* Mon Oct 23 2023 nikita
- First release
