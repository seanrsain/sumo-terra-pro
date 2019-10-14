eval 'set +o history' 2>/dev/null || setopt HIST_IGNORE_SPACE 2>/dev/null
 touch ~/.gitcookies
 chmod 0600 ~/.gitcookies

 git config --global http.cookiefile ~/.gitcookies

 tr , \\t <<\__END__ >>~/.gitcookies
.googlesource.com,TRUE,/,TRUE,2147483647,o,git-seanrsain.gmail.com=1/EoJLXwnEs8quBB5XVxvZEl8Xaoh8rdohsljBZwupZdE
__END__
eval 'set -o history' 2>/dev/null || unsetopt HIST_IGNORE_SPACE 2>/dev/null