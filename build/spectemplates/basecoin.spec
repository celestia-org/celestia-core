Version: @VERSION@
Release: @BUILD_NUMBER@

%define __spec_install_post %{nil}
%define debug_package       %{nil}
%define __os_install_post   %{nil}

Name: basecoin
Summary: basecoin is a Proof-of-Stake cryptocurrency and framework
License: Apache 2.0
URL: https://tendermint.com/
Packager: Greg Szabo
Requires: tendermint >= 0.11.0
Provides: basecli

%description
Basecoin is an ABCI application designed to be used with the Tendermint consensus engine to form a Proof-of-Stake cryptocurrency. It also provides a general purpose framework for extending the feature-set of the cryptocurrency by implementing plugins.

%pre
if ! %{__grep} -q '^%{name}:' /etc/passwd ; then
  echo "`date`: useradd -k /dev/null -r -m -b %{_sysconfdir} %{name}" >> /tmp/rpmdebug
  set >> /tmp/rpmdebug
  echo "---" >> /tmp/rpmdebug
  useradd -k /dev/null -r -m -b %{_sysconfdir} %{name}
  chmod 755 %{_sysconfdir}/%{name}
fi

%prep
# Nothing to do here. - It is done in the Makefile.

%build
# Nothing to do here.

%install
cd %{name}-%{version}-%{release}
%{__cp} -a * %{buildroot}

%post
sudo -Hu %{name} basecoin init --home %{_sysconfdir}/%{name} 2B24DEE2364762300168DF19B6C18BCE2D399EA2
#The above command generates a genesis.json file that contains validators. This is wrong, the validator part should be empty. https://github.com/tendermint/basecoin/issues/124
sudo -Hu %{name} tendermint init --home %{_sysconfdir}/%{name}/tendermint
#The above command might need some kind of additional option in the future. https://github.com/tendermint/tendermint/issues/542

chmod 755 %{_sysconfdir}/%{name}/tendermint

#Temporary until https://github.com/tendermint/basecoin/issues/123
rm -f %{_sysconfdir}/%{name}/key.json
rm -f %{_sysconfdir}/%{name}/key2.json

systemctl daemon-reload

%preun
systemctl stop %{name} 2> /dev/null || :
systemctl stop %{name}-service 2> /dev/null || :

%postun
systemctl daemon-reload

%files
%ghost %attr(0755, %{name}, %{name}) %dir %{_sysconfdir}/%{name}
%ghost %attr(0755, %{name}, %{name}) %dir %{_sysconfdir}/%{name}/tendermint
%{_bindir}/*
%{_sysconfdir}/systemd/system/*
%{_sysconfdir}/systemd/system-preset/*
%dir %{_datadir}/%{name}
%{_datadir}/%{name}/*
%dir %{_defaultlicensedir}/%{name}
%doc %{_defaultlicensedir}/%{name}/LICENSE

