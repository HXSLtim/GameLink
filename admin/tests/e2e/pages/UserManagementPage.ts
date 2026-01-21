import { Page } from '@playwright/test';

/**
 * Page Object Model for User Management Page
 * Encapsulates all user management interactions and assertions
 */
export class UserManagementPage {
  constructor(private page: Page) {}

  // Element locators - PageContainer uses h1 for title
  private readonly pageTitle = this.page.locator('h1').filter({ hasText: /用户管理/ });
  private readonly searchInput = this.page.getByPlaceholder(/搜索|search/i);
  private readonly createButton = this.page.getByRole('button', { name: /新增|创建|create|add/i });
  private readonly table = this.page.locator('.ant-table');
  private readonly tableRows = this.page.locator('.ant-table-tbody tr');
  private readonly editButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /编辑|edit/i });
  private readonly deleteButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /删除|delete/i });
  private readonly confirmDeleteButton = this.page.getByRole('button', { name: /确认|confirm/i });
  private readonly cancelButton = this.page.getByRole('button', { name: /取消|cancel/i });

  // Modal locators
  private readonly modal = this.page.locator('.ant-modal');
  private readonly modalTitle = this.modal.locator('.ant-modal-title');
  private readonly modalOkButton = this.modal.getByRole('button', { name: /确定|ok|submit/i });
  private readonly modalCancelButton = this.modal.getByRole('button', { name: /取消|cancel/i });

  // Form fields in modal
  private readonly nameInput = this.modal.getByLabel(/姓名|name/i);
  private readonly emailInput = this.modal.getByLabel(/邮箱|email/i);
  private readonly phoneInput = this.modal.getByLabel(/电话|phone/i);
  private readonly roleSelect = this.modal.getByLabel(/角色|role/i);
  private readonly statusSelect = this.modal.getByLabel(/状态|status/i);

  /**
   * Navigate to user management page
   */
  async goto() {
    await this.page.goto('/admin/sys/user');
    await this.waitForPageLoad();
  }

  /**
   * Wait for page to load completely
   */
  async waitForPageLoad() {
    await this.page.waitForLoadState('networkidle');
    // Wait for either the page title or the table to be visible
    await Promise.race([
      this.pageTitle.waitFor({ state: 'visible', timeout: 15000 }),
      this.table.waitFor({ state: 'visible', timeout: 15000 }),
    ]);
  }

  /**
   * Search for a user by keyword
   */
  async searchUser(keyword: string) {
    await this.searchInput.fill(keyword);
    await this.page.waitForTimeout(500); // Wait for debounce
  }

  /**
   * Click create user button
   */
  async clickCreateUser() {
    await this.createButton.first().click();
    await expect(this.modal).toBeVisible();
  }

  /**
   * Fill user form in modal
   */
  async fillUserForm(userData: {
    name: string;
    email: string;
    phone: string;
    role: string;
    status?: string;
    password?: string;
  }) {
    await this.nameInput.fill(userData.name);
    await this.emailInput.fill(userData.email);
    await this.phoneInput.fill(userData.phone);

    // Handle role selection
    await this.roleSelect.click();
    await this.page.getByRole('option', { name: userData.role }).click();

    // Handle status if provided
    if (userData.status) {
      await this.statusSelect.click();
      await this.page.getByRole('option', { name: userData.status }).click();
    }

    // Handle password if provided
    if (userData.password) {
      const passwordInput = this.modal.getByLabel(/密码|password/i);
      await passwordInput.fill(userData.password);
    }
  }

  /**
   * Submit user form
   */
  async submitForm() {
    await this.modalOkButton.click();
    await this.page.waitForTimeout(1000); // Wait for API call
  }

  /**
   * Create a new user
   */
  async createUser(userData: {
    name: string;
    email: string;
    phone: string;
    role: string;
    status?: string;
    password?: string;
  }) {
    await this.clickCreateUser();
    await this.fillUserForm(userData);
    await this.submitForm();
  }

  /**
   * Edit a user from table row
   */
  async editUser(rowIndex: number, updatedData: Partial<{
    name: string;
    email: string;
    phone: string;
    role: string;
    status: string;
  }>) {
    await this.editButton(rowIndex).click();
    await expect(this.modal).toBeVisible();

    if (updatedData.name) await this.nameInput.fill(updatedData.name);
    if (updatedData.email) await this.emailInput.fill(updatedData.email);
    if (updatedData.phone) await this.phoneInput.fill(updatedData.phone);
    if (updatedData.role) {
      await this.roleSelect.click();
      await this.page.getByRole('option', { name: updatedData.role }).click();
    }
    if (updatedData.status) {
      await this.statusSelect.click();
      await this.page.getByRole('option', { name: updatedData.status }).click();
    }

    await this.submitForm();
  }

  /**
   * Delete a user from table row
   */
  async deleteUser(rowIndex: number) {
    await this.deleteButton(rowIndex).click();
    await expect(this.modal).toBeVisible();
    await this.confirmDeleteButton.click();
    await this.page.waitForTimeout(1000);
  }

  /**
   * Cancel delete operation
   */
  async cancelDelete(rowIndex: number) {
    await this.deleteButton(rowIndex).click();
    await this.modalCancelButton.click();
  }

  /**
   * Get user count from table
   */
  async getUserCount(): Promise<number> {
    await this.page.waitForLoadState('networkidle');
    return await this.tableRows.count();
  }

  /**
   * Get text from table cell
   */
  async getCellText(rowIndex: number, columnIndex: number): Promise<string> {
    const cell = this.tableRows.nth(rowIndex).locator('td').nth(columnIndex);
    return await cell.textContent() || '';
  }

  /**
   * Verify user exists in table
   */
  async verifyUserExists(name: string): Promise<boolean> {
    const userCount = await this.getUserCount();
    for (let i = 0; i < userCount; i++) {
      const cellText = await this.getCellText(i, 1); // Assuming name is in column 1
      if (cellText.includes(name)) {
        return true;
      }
    }
    return false;
  }

  /**
   * Wait for user to appear in table
   */
  async waitForUserToAppear(name: string, timeout = 10000) {
    await this.page.waitForTimeout(1000); // Initial wait
    const startTime = Date.now();

    while (Date.now() - startTime < timeout) {
      if (await this.verifyUserExists(name)) {
        return true;
      }
      await this.page.waitForTimeout(500);
    }

    throw new Error(`User "${name}" did not appear in table within ${timeout}ms`);
  }

  /**
   * Verify success message
   */
  async verifySuccessMessage() {
    await expect(this.page.locator('.ant-message-success')).toBeVisible();
  }

  /**
   * Verify error message
   */
  async verifyErrorMessage() {
    await expect(this.page.locator('.ant-message-error')).toBeVisible();
  }

  /**
   * Navigate to next page if pagination exists
   */
  async nextPage() {
    const nextButton = this.page.getByRole('button', { name: /下一页|next/i });
    if (await nextButton.isEnabled()) {
      await nextButton.click();
      await this.page.waitForLoadState('networkidle');
    }
  }

  /**
   * Filter by role
   */
  async filterByRole(role: string) {
    const roleFilter = this.page.getByText(/角色|role/i).first();
    await roleFilter.click();
    await this.page.getByRole('option', { name: role }).click();
    await this.page.waitForTimeout(1000);
  }

  /**
   * Filter by status
   */
  async filterByStatus(status: string) {
    const statusFilter = this.page.getByText(/状态|status/i).first();
    await statusFilter.click();
    await this.page.getByRole('option', { name: status }).click();
    await this.page.waitForTimeout(1000);
  }

  /**
   * Click user row to view details
   */
  async viewUserDetails(rowIndex: number) {
    await this.tableRows.nth(rowIndex).click();
    await this.page.waitForLoadState('networkidle');
  }

  /**
   * Batch delete users
   */
  async batchDeleteUsers(userNames: string[]) {
    // Select checkboxes
    for (const name of userNames) {
      const checkbox = this.page.getByText(name).locator('..').getByRole('checkbox');
      await checkbox.check();
    }

    // Click batch delete button
    const batchDeleteButton = this.page.getByRole('button', { name: /批量删除/i });
    await batchDeleteButton.click();

    // Confirm deletion
    await this.confirmDeleteButton.click();
    await this.page.waitForTimeout(1000);
  }

  /**
   * Export user list
   */
  async exportUsers() {
    const exportButton = this.page.getByRole('button', { name: /导出|export/i });
    await exportButton.click();
    await this.page.waitForTimeout(2000);
  }
}
