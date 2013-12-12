using GalaSoft.MvvmLight.Command;
using System.Windows.Input;

namespace famousfront.core
{
  abstract class DialogViewModel : ViewModelBase
  {
    [System.Flags]
    internal enum CommandFlags
    {
      HasNone = 0,
      HasClose = 1,
      HasCancel = 2,
      HasConfirm = 4,
    }
    protected DialogViewModel(CommandFlags flag)
    {
      HasCancel = flag.HasFlag(CommandFlags.HasCancel);
      HasClose = flag.HasFlag(CommandFlags.HasClose);
      HasConfirm = flag.HasFlag(CommandFlags.HasConfirm);
    }
    bool _has_cancel;
    bool _has_close;
    bool _has_confirm;
    public bool HasCancel
    {
      get { return _has_cancel; }
      set { Set(ref _has_cancel, value); }
    }
    public bool HasClose
    {
      get { return _has_close; }
      set { Set(ref _has_close, value); }
    }
    public bool HasConfirm
    {
      get { return _has_confirm; }
      set { Set(ref _has_confirm, value); }
    }
    ICommand _confirm_command;
    public ICommand ConfirmCommand
    {
      get {return _confirm_command ?? (_confirm_command = confirm_command());        }
    }
    ICommand _cancel_command;
    public ICommand CancelCommand
    {
      get { return _cancel_command ?? (_cancel_command = cancel_command()); }
    }
    ICommand _close_command;
    public ICommand CloseCommand
    {
      get { return _close_command ?? (_close_command = close_command()); }
    }
    ICommand confirm_command()
    {
      return new RelayCommand(ExecuteConfirm, CanExecuteConfirm);
    }

    protected virtual bool CanExecuteConfirm()
    {
      return true;
    }

    protected virtual void ExecuteConfirm()
    {
    }
    ICommand cancel_command()
    {
      return new RelayCommand(ExecuteCancel, CanExecuteCancel);
    }

    protected virtual bool CanExecuteCancel()
    {
      return true;
    }

    protected virtual void ExecuteCancel()
    {
    }
    ICommand close_command()
    {
      return new RelayCommand(ExecuteClose, CanExecuteClose);
    }

    protected virtual bool CanExecuteClose()
    {
      return true;
    }

    protected virtual void ExecuteClose()
    {
    }
  }
  class OverlappedViewModel : DialogViewModel
  {
    internal OverlappedViewModel() : base(CommandFlags.HasClose)
    {

    }
  }
}
