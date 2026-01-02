$pass = ConvertTo-SecureString -String 'Vagrant123!' -AsPlainText -Force
New-LocalUser -Name 'vagrant2' -Password $pass -AccountNeverExpires
Add-LocalGroupMember -Group 'Administrators' -Member 'vagrant2'
Add-LocalGroupMember -Group 'Users' -Member 'vagrant2'
