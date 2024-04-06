import { Module } from '@nestjs/common';
import { RegisterModule } from './register/register.module';
import { LoginModule } from './login/login.module';
import { LogoutModule } from './logout/logout.module';
import { VerifyCodeModule } from './verify-code/verify-code.module';
import { RefreshModule } from './refresh/refresh.module';

@Module({
  imports: [RegisterModule, LoginModule, LogoutModule, VerifyCodeModule, RefreshModule],
  controllers: [],
})
export class AuthModule {}
