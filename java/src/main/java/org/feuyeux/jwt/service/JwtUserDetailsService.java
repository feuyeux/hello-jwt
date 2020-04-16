package org.feuyeux.jwt.service;

import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.core.userdetails.User;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.stereotype.Service;

import java.util.ArrayList;

@Service
public class JwtUserDetailsService implements UserDetailsService {

    public static final SimpleGrantedAuthority READ_ROLE = new SimpleGrantedAuthority("READ_ROLE");

    @Override
    public UserDetails loadUserByUsername(String username) throws UsernameNotFoundException {
        if ("hello_man".equals(username)) {
            ArrayList<GrantedAuthority> authorities = new ArrayList<>();
            authorities.add(READ_ROLE);
            authorities.add(new SimpleGrantedAuthority("WRITE_ROLE"));
            return new User("hello_man", "$2a$10$slYQmyNdGzTn7ZLBXBChFOC9f6kFjAqPhccnP6DxlWXx2lPk1C3G6",
                    authorities);
        } else {
            throw new UsernameNotFoundException("User not found with username: " + username);
        }
    }
}