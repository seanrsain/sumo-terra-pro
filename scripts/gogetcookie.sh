touch ~/.gitcookies
chmod 0600 ~/.gitcookies

git config --global http.cookiefile ~/.gitcookies

tr , \\t <<\__END__ >>~/.gitcookies
.googlesource.com,TRUE,/,TRUE,2147483647,o,git-seanrsain.gmail.com=1/EoJLXwnEs8quBB5XVxvZEl8Xaoh8rdohsljBZwupZdE
__END__